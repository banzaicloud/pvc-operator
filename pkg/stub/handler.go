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
		if o.Status.Phase == v1.ClaimPending {
			var commonProvider providers.CommonProvider
			logrus.Info("PersistenVolumeClaim event received!")
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
	return nil
}

//func GetAccountKeys(ctx context.Context, accountName string, subscriptionID string, resourceGroupName string) (storage.AccountListKeysResult, error) {
//	accountsClient, _ := createStorageAccountClient(subscriptionID)
//	return accountsClient.ListKeys(ctx, resourceGroupName, accountName)
//}

//func regenerateAccountKey(ctx context.Context, accountName string, key int, subscriptionID string, resourceGroupName string) (list storage.AccountListKeysResult, err error) {
//	oldKeys, err := GetAccountKeys(ctx, accountName, subscriptionID, resourceGroupName)
//	if err != nil {
//		return list, err
//	}
//	accountsClient, err := createStorageAccountClient(subscriptionID)
//	return accountsClient.RegenerateKey(
//		ctx,
//		resourceGroupName,
//		accountName,
//		storage.AccountRegenerateKeyParameters{
//			KeyName: (*oldKeys.Keys)[key].KeyName,
//		})
//}
