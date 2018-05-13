package providers

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/coreos/operator-sdk/pkg/sdk/query"
	"fmt"
	"k8s.io/api/apps/v1beta1"
	"github.com/coreos/operator-sdk/pkg/sdk/action"
)

type CommonProvider interface {
	CreateStorageClass(*v1.PersistentVolumeClaim) error
	GenerateMetadata() error
	determineParameters(*v1.PersistentVolumeClaim) (map[string]string, error)
	determineProvisioner(*v1.PersistentVolumeClaim) (string, error)
}

func DetermineProvider() (CommonProvider, error) {
	var providers = map[string]string{
		"azure": "http://169.254.169.254/metadata/instance?api-version=2017-12-01",
		"aws":   "http://169.254.169.254/latest/meta-data/",
		"google": "http://169.254.169.254/0.1/meta-data/",
	}
	for key, value := range providers {
		req, err := http.NewRequest("GET", value, nil)
		if err != nil {
			logrus.Errorf("Could not create a proper http request %s", err.Error())
			return nil, err
		}
		req.Header.Set("Metadata", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Something happened during the request %s", err.Error())
		}
		if resp.StatusCode == 404 || resp.StatusCode == 405 {
			continue
		}
		switch key {
		case "azure":
			return &AzureProvider{}, nil
		case "aws":
			return &AwsProvider{}, nil
		case "google":
			return &GoogleProvider{}, nil
		}

	}
	return nil, fmt.Errorf("could not determine cloud provider")
}

func CheckNfsServerExistence(name string) bool {
	if !CheckPersistentVolumeClaimExistence("nfs") {
		logrus.Info("PersistentVolume claim for Nfs does not exists!")
		return false
	}
	if !checkNfsProviderDeployment() {
		logrus.Info("Nfs provider deployment does not exists!")
		return false
	}
	if !CheckStorageClassExistence(name) {
		logrus.Info("StorageClass for Nfs does not exist!")
		return false
	}
	return true

}

func CheckPersistentVolumeClaimExistence(name string) bool {
	persistentVolumeClaim := &v1.PersistentVolumeClaim{
		TypeMeta: v12.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name: name,
		},
	}
	if err := query.Get(persistentVolumeClaim); err != nil {
		logrus.Info("PersistentVolumeClaim does not exists %s", err.Error())
		return false
	}
	logrus.Infof("PersistentVolumeClaim %s exist!", name)
	return true
}

func checkNfsProviderDeployment() bool {
	deployment := &v1beta1.Deployment{
		TypeMeta: v12.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1beta1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:            "nfs-provisioner",
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers:                    []v1.Container{
						{
							Args: []string{
								"-provisioner=banzaicloud.com/nfs",
							},
						},
					},
				},
			},
		},
	}
	service := &v1.Service{
		TypeMeta: v12.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:            "nfs-provisioner",
		},
		Spec: v1.ServiceSpec{
			Selector:                 map[string]string{"app": "nfs-provisioner"},
		},
	}
	if err := query.Get(deployment); err != nil {
		logrus.Infof("Nfs provider deployment does not exists %s", err.Error())
		return false
	}
	if err := query.Get(service); err != nil {
		logrus.Infof("Nfs provider service does not exists %s", err.Error())
		return false
	}
	logrus.Info("Nfs provider exists!")
	return true
}

func SetUpNfsProvisioner(pv *v1.PersistentVolumeClaim) error {
	logrus.Info("Creating new PersistentVolumeClaim for Nfs provisioner..")
	err := action.Create(&v1.PersistentVolumeClaim{
		TypeMeta: v12.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:            "nfs-data",
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": *pv.Spec.Resources.Requests.StorageEphemeral(),
				},
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error happened during creating a PersistentVolumeClaim for Nfs %s",err.Error())
		return err
	}
	logrus.Info("Creating new Service for Nfs provisioner..")
	err = action.Create(&v1.Service{
		TypeMeta: v12.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:            "nfs-provisioner",
			Labels:                     map[string]string{
				"app": "nfs-provisioner",
			},
		},
		Spec: v1.ServiceSpec{
			Ports:                    []v1.ServicePort{
				{Name: "nfs", Port: 2049},
				{Name: "mountd", Port: 20048},
				{Name: "rpcbind", Port: 111},
				{Name: "rpcbind-udp", Port: 111, Protocol: "UDP"},
			},
			Selector:                 map[string]string{
				"app": "nfs-provisioner",
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error happened during creating the Service for Nfs %s", err.Error())
		return err
	}
	logrus.Info("Creating new Deployment for Nfs provisioner..")
	replicas := int32(1)
	action.Create(&v1beta1.Deployment{
		TypeMeta: v12.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1beta1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:            "nfs-provisioner",
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Labels:                     map[string]string{
						"app": "nfs-provisioner",
					},
				},
				Spec: v1.PodSpec{
					Volumes:                       []v1.Volume{
						{Name: "nfs-prov-volume", VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: "nfs-data",
							},
						},},
					},
					Containers:                    []v1.Container{
						{},
					},
				},
			},
			Strategy: v1beta1.DeploymentStrategy{
				Type: "Recreate",

			},
		},
	})
	return nil
}
func CheckStorageClassExistence(name string) bool {
	storageClass := &storagev1.StorageClass{
		TypeMeta: v12.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name: name,
		},
	}
	if err := query.Get(storageClass); err != nil {
		logrus.Infof("Storageclass does not exist %s", err.Error())
		return false
	}
	logrus.Infof("Storageclass %s exists!", name)
	return true
}
