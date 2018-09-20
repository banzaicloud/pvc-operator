package stub

import (
	"fmt"
	"github.com/banzaicloud/pvc-operator/pkg/apis/banzaicloud/v1alpha1"
	"github.com/banzaicloud/pvc-operator/pkg/stub/providers"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"strings"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx sdk.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1.PersistentVolumeClaim:
		if o.Spec.StorageClassName != nil {
			if o.Status.Phase == v1.ClaimPending {
				logrus.Info("PersistenVolumeClaim event received!")
				logrus.Info("Check if the storageclass already exist!")
				if strings.Contains(*o.Spec.StorageClassName, "nfs") {
					logrus.Info("Check if the deployment for Nfs exists!")
					if !providers.CheckNfsServerExistence(*o.Spec.StorageClassName, o.Namespace) {
						err := providers.SetUpNfsProvisioner(o)
						if err != nil {
							logrus.Errorf("Cloud not create the NFS deployment %s", err.Error())
							return err
						}
					}
					return nil
				}
				if !providers.CheckStorageClassExistence(*o.Spec.StorageClassName) {
					commonProvider, err := providers.DetermineProvider()
					if err != nil {
						logrus.Errorf("Cloud not determine cloud provider %s", err.Error())
						return err
					}
					if err := commonProvider.GenerateMetadata(); err != nil {
						logrus.Errorf("Cloud not generate metadata %s", err.Error())
						return err
					}
					if err := commonProvider.CreateStorageClass(o); err != nil && !apierrors.IsAlreadyExists(err) {
						logrus.Errorf("Failed to create a storageclass: %s", err.Error())
						return fmt.Errorf("failed to create storageclass: %s", err.Error())
					}
					return nil
				}
			}
		}
	case *v1alpha1.ObjectStore:
		logrus.Info("Object Store creation event received!")
		logrus.Info("Check of the bucket already exists!")
		commonProvider, err := providers.DetermineProvider()
		if err != nil {
			logrus.Errorf("Cloud not determine cloud provider %s", err.Error())
			return err
		}
		if err := commonProvider.CreateObjectStoreBucket(o); err != nil {
			logrus.Errorf("Could not create an ObjectStore Bucket %s", err.Error())
		}
		return nil
	}
	return nil
}
