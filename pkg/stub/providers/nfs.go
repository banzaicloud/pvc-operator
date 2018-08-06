package providers

import (
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetUpNfsProvisioner sets up a deployment a pvc and a service to handle nfs workload
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
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": parsedStorageSize,
				},
			},
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Error happened during creating a PersistentVolumeClaim for Nfs %s", err.Error())
		return err
	}
	logrus.Info("Creating new Service for Nfs provisioner..")
	err = sdk.Create(&v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nfs-provisioner",
			Labels: map[string]string{
				"app": "nfs-provisioner",
			},
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{Name: "nfs", Port: 2049},
				{Name: "mountd", Port: 20048},
				{Name: "rpcbind", Port: 111},
				{Name: "rpcbind-udp", Port: 111, Protocol: "UDP"},
			},
			Selector: map[string]string{
				"app": "nfs-provisioner",
			},
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Error happened during creating the Service for Nfs %s", err.Error())
		return err
	}
	logrus.Info("Creating new Deployment for Nfs provisioner..")
	replicas := int32(1)
	err = sdk.Create(&v1beta1.Deployment{
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
					Labels: map[string]string{
						"app": "nfs-provisioner",
					},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "nfs-prov-volume", VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: fmt.Sprintf("%s-data", *pv.Spec.StorageClassName),
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  "nfs-provisioner",
							Image: "quay.io/kubernetes_incubator/nfs-provisioner:v1.0.9",
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
								{Name: "POD_IP", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "status.podIP"}}},
								{Name: "SERVICE_NAME", Value: "nfs-provisioner"},
								{Name: "POD_NAMESPACE", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.namespace"}}},
							},
							VolumeMounts: []v1.VolumeMount{
								{Name: "nfs-prov-volume", MountPath: "/export"},
							},
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{v1.ResourceCPU: resource.MustParse("1")},
							},
						},
					},
				},
			},
			Strategy: v1beta1.DeploymentStrategy{
				Type: v1beta1.RecreateDeploymentStrategyType,
			},
		},
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Error happened during creating the Deployment for Nfs %s", err.Error())
		return err
	}
	logrus.Info("Creating new StorageClass for Nfs provisioner..")
	reclaimPolicy := v1.PersistentVolumeReclaimRetain
	err = sdk.Create(&storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: *pv.Spec.StorageClassName,
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
		ReclaimPolicy: &reclaimPolicy,
		Provisioner:   "banzaicloud.com/nfs",
	})
	if err != nil && !errors.IsAlreadyExists(err) {
		logrus.Errorf("Error happened during creating the StorageClass for Nfs %s", err.Error())
		return err
	}
	return nil
}

// CheckNfsServerExistence checks if the NFS deployment and all companion service exists
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

// checkNfsProviderDeployment checks if the NFS deployment exists
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
					Containers: []v1.Container{
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
			Selector: map[string]string{"app": "nfs-provisioner"},
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
