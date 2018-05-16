package providers

import (
	"k8s.io/api/core/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"fmt"

	"k8s.io/api/apps/v1beta1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func SetUpNfsProvisioner(pv *v1.PersistentVolumeClaim) error {
	logrus.Info("Creating new PersistentVolumeClaim for Nfs provisioner..")

	ownerRef := asOwner(getOwner())
	plusStorage, _ := resource.ParseQuantity("2Gi")
	parsedStorageSize := pv.Spec.Resources.Requests["storage"]
	parsedStorageSize.Add(plusStorage)

	err := sdk.Create(&v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-data", *pv.Spec.StorageClassName),
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": parsedStorageSize,
				},
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error happened during creating a PersistentVolumeClaim for Nfs %s",err.Error())
		return err
	}
	logrus.Info("Creating new Service for Nfs provisioner..")
	err = sdk.Create(&v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "nfs-provisioner",
			Labels:                     map[string]string{
				"app": "nfs-provisioner",
			},
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
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
	sdk.Create(&v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs-provisioner",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:                     map[string]string{
						"app": "nfs-provisioner",
					},
				},
				Spec: v1.PodSpec{
					Volumes:                       []v1.Volume{
						{Name: "nfs-prov-volume", VolumeSource: v1.VolumeSource{
							PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
								ClaimName: fmt.Sprintf("%s-data", *pv.Spec.StorageClassName),
							},
						},},
					},
					Containers:                    []v1.Container{
						{
							Name: "nfs-provisioner",
							Image: "quay.io/kubernetes_incubator/nfs-provisioner:v1.0.8",
							Ports: []v1.ContainerPort{
								{Name: "nfs", ContainerPort: 2049},
								{Name: "mountd", ContainerPort: 20048},
								{Name: "rpcbind", ContainerPort: 111},
								{Name: "rpcbind-udp", ContainerPort: 111, Protocol: "UDP"},
							},
							SecurityContext: &v1.SecurityContext{
								Capabilities: &v1.Capabilities{
									Add: []v1.Capability{
										"DAC_READ_SEARCH",
										"SYS_RESOURCE",
									},
								},
							},
							Args: []string{
								"-provisioner=banzaicloud.com/nfs",
							},
							Env: []v1.EnvVar{
								{Name: "POD_IP", ValueFrom: &v1.EnvVarSource{ FieldRef: &v1.ObjectFieldSelector{FieldPath: "status.podIP"}}},
								{Name: "SERVICE_NAME", Value: "nfs-provisioner"},
								{Name: "POD_NAMESPACE", ValueFrom: &v1.EnvVarSource{ FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.namespace"}}},
							},
							VolumeMounts: []v1.VolumeMount{
								{Name:"nfs-prov-volume", MountPath: "/export"},
							},
						},
					},
				},
			},
			Strategy: v1beta1.DeploymentStrategy{
				Type: "Recreate",

			},
		},
	})
	if err != nil {
		logrus.Errorf("Error happened during creating the Deployment for Nfs %s", err.Error())
		return err
	}
	logrus.Info("Creating new StorageClass for Nfs provisioner..")
	sdk.Create(&storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            *pv.Spec.StorageClassName,
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		Provisioner:          "banzaicloud.com/nfs",
	})
	if err != nil {
		logrus.Errorf("Error happened during creating the StorageClass for Nfs %s", err.Error())
		return err
	}
	return nil
}

func CheckNfsServerExistence(name string) bool {
	if !CheckPersistentVolumeClaimExistence(fmt.Sprintf("%s-data", name)) {
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

func checkNfsProviderDeployment() bool {
	deployment := &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs-provisioner",
			Namespace: "default",
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
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs-provisioner",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Selector:                 map[string]string{"app": "nfs-provisioner"},
		},
	}
	if err := sdk.Get(deployment); err != nil {
		logrus.Infof("Nfs provider deployment does not exists %s", err.Error())
		return false
	}
	if err := sdk.Get(service); err != nil {
		logrus.Infof("Nfs provider service does not exists %s", err.Error())
		return false
	}
	logrus.Info("Nfs provider exists!")
	return true
}
