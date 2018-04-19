package stub

import (
	"github.com/coreos/operator-sdk/pkg/sdk/handler"
	"github.com/coreos/operator-sdk/pkg/sdk/types"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/api/core/v1"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/banzaicloud/pvc-handler/pkg/stub/providers"
)

func NewHandler() handler.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx types.Context, event types.Event) error {
	switch o := event.Object.(type) {
	case *v1.PersistentVolumeClaim:
		if o.Spec.StorageClassName != nil {
			if o.Status.Phase == v1.ClaimPending {
				logrus.Info("PersistenVolumeClaim event received!")
				logrus.Info("Check if the storageclass already exist!")
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
	}
	return nil
}