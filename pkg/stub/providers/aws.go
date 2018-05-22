package providers

import (
	"errors"
	"github.com/banzaicloud/pvc-operator/pkg/apis/banzaicloud/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AwsProvider struct {
}

func (aws *AwsProvider) CreateStorageClass(pvc *v1.PersistentVolumeClaim) error {
	logrus.Info("Creating new storage class")
	provisioner, err := aws.determineProvisioner(pvc)
	if err != nil {
		return err
	}
	logrus.Info("Determining provisioner succeeded")
	parameter, err := aws.determineParameters(pvc)
	if err != nil {
		return err
	}
	logrus.Info("Determining parameter succeeded")
	return sdk.Create(&storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            *pvc.Spec.StorageClassName,
			Annotations:     nil,
			OwnerReferences: nil,
		},
		Provisioner:  provisioner,
		MountOptions: nil,
		Parameters:   parameter,
	})
}

func (aws *AwsProvider) GenerateMetadata() error {
	return nil
}

func (aws *AwsProvider) determineParameters(pvc *v1.PersistentVolumeClaim) (map[string]string, error) {
	//var parameter = map[string]string{}
	for _, mode := range pvc.Spec.AccessModes {
		switch mode {
		case "ReadWriteOnce":
			return nil, nil
		}
	}
	return nil, errors.New("could not determine parameters")
}

func (aws *AwsProvider) determineProvisioner(pvc *v1.PersistentVolumeClaim) (string, error) {
	for _, mode := range pvc.Spec.AccessModes {
		switch mode {
		case "ReadWriteOnce":
			return "kubernetes.io/aws-ebs", nil
		case "ReadWriteMany":
			return "", errors.New("Not supported yet")
		case "ReadOnlyMany":
			return "", errors.New("Not supported yet")
		}
	}
	return "", errors.New("AccessMode is missing from the PVC")
}

func (aws *AwsProvider) CheckBucketExistence(store *v1alpha1.ObjectStore) (bool, error) {
	return false, nil
}

func (aws *AwsProvider) CreateObjectStoreBucket(store *v1alpha1.ObjectStore) error {
	return nil
}
