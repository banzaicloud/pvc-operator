package providers

import (
	"k8s.io/api/core/v1"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/coreos/operator-sdk/pkg/sdk/action"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GoogleProvider struct {
}

func (gke *GoogleProvider) CreateStorageClass(pvc *v1.PersistentVolumeClaim) error {
	logrus.Info("Creating new storage class")
	provisioner, err := gke.determineProvisioner(pvc)
	if err != nil {
		return err
	}
	logrus.Info("Determining provisioner succeeded")
	parameter, err :=  gke.determineParameters(pvc)
	if err != nil {
		return err
	}
	logrus.Info("Determining parameter succeeded")
	return action.Create(&storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            *pvc.Spec.StorageClassName,
			Annotations:                nil,
			OwnerReferences:            nil,
		},
		Provisioner:          provisioner,
		MountOptions:         nil,
		Parameters:           parameter,
	})
}

func (gke *GoogleProvider) GenerateMetadata() error {
	return nil
}

func (gke *GoogleProvider) determineParameters(pvc *v1.PersistentVolumeClaim) (map[string]string, error) {
	//var parameter = map[string]string{}
	for _, mode := range pvc.Spec.AccessModes {
		switch mode {
		case "ReadWriteOnce", "ReadOnlyMany":
			return nil, nil
		}
	}
	return nil, errors.New("could not determine parameters")
}

func (gke *GoogleProvider) determineProvisioner (pvc *v1.PersistentVolumeClaim) (string, error) {
	for _, mode := range pvc.Spec.AccessModes {
		switch mode {
		case "ReadWriteOnce", "ReadOnlyMany":
			return "kubernetes.io/gce-pd", nil
		case "ReadWriteMany":
			return "", errors.New("Not supported yet")
		}
	}
	return "", errors.New("AccessMode is missing from the PVC")
}