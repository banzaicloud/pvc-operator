package providers

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
	"github.com/coreos/operator-sdk/pkg/sdk/action"
)

type Metadata struct {
	location string
	subscriptionID string
	resourceGroupName string
}

type AzureProvider struct {
	metadata Metadata
}

const (
	storageAccount = "storageAccount"
	location = "location"
	skuName = "skuName"
	kind = "kind"

)
func (az *AzureProvider) CreateStorageClass (pvc *v1.PersistentVolumeClaim) error {
	logrus.Info("Creating new storage class")
	provisioner, err := determineProvisioner(pvc)
	if err != nil {
		return nil
	}
	logrus.Info("Determining provisioner succeeded")
	parameter, err :=  determineParameters(pvc, az)
	if err != nil {
		return nil
	}
	logrus.Info("Determining parameter succeeded")
	if parameter[storageAccount] != "" {
		createStorageAccount(context.TODO(), parameter[storageAccount], az)
	}
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

func (az *AzureProvider) GenerateMetadata() error {
	logrus.Infof("Getting Metadata from service")
	var metadatas = [3]string{
		"location",
		"subscriptionId",
		"resourceGroupName",
	}
	var result = map[string]string{}
	for _, metadata := range metadatas{
		req, err := http.NewRequest("GET", fmt.Sprintf("http://169.254.169.254/metadata/instance/compute/%s?api-version=2017-12-01&format=text", metadata), nil)
		if err != nil {
			logrus.Errorf("Error during getting %s, %s", metadata, err.Error())
			return err
		}
		req.Header.Add("Metadata", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logrus.Errorf("Error during getting %s, %s", metadata, err.Error())
			return err
		}
		readMetadata, err :=ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Error during reading response %s", err.Error())
			return err
		}
		result[metadata] = string(readMetadata)
	}
	az.metadata.location = result["location"]
	az.metadata.subscriptionID = result["subscriptionId"]
	az.metadata.resourceGroupName = result["resourceGroupName"]

	return nil
}

func createStorageAccount(ctx context.Context, accountName string, az *AzureProvider) (s storage.Account, err error) {
	storageAccountsClient, err := createStorageAccountClient(az.metadata.subscriptionID)

	result, err := storageAccountsClient.CheckNameAvailability(
		ctx,
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(accountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if err != nil {
		logrus.Fatalf("%s: %v", "storage account creation failed", err)
	}
	if *result.NameAvailable != true {
		logrus.Fatalf("%s [%s]: %v: %v", "storage account name not available", accountName, err, *result.Message)
	}

	future, err := storageAccountsClient.Create(
		ctx,
		az.metadata.resourceGroupName,
		accountName,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS},
			Kind:     storage.Storage,
			Location: to.StringPtr(az.metadata.location),
			AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{},
		})

	if err != nil {
		return s, fmt.Errorf("cannot create storage account: %v", err)
	}

	err = future.WaitForCompletion(ctx, storageAccountsClient.Client)
	if err != nil {
		return s, fmt.Errorf("cannot get the storage account create future response: %v", err)
	}
	logrus.Info("StorageAccount created!")
	return future.Result(storageAccountsClient)
}

func createStorageAccountClient(subscriptionID string) (storage.AccountsClient, error) {
	accountClient := storage.NewAccountsClient(subscriptionID)
	logrus.Info("Authenticating...")
	authorizer, err := auth.NewMSIConfig().Authorizer()
	if err != nil {
		logrus.Errorf("Error happened during authentication %s", err.Error())
		return storage.AccountsClient{}, err
	}
	accountClient.Authorizer = authorizer
	logrus.Info("Authenticating succeeded")
	return accountClient, nil
}

func determineParameters(pvc *v1.PersistentVolumeClaim, provider *AzureProvider) (map[string]string, error) {
	var parameter = map[string]string{}
	for _, mode := range pvc.Spec.AccessModes {
		switch mode {
		case "ReadWriteOnce":
			parameter[skuName] = "Standard_LRS"
			parameter[kind] = "managed"
			return parameter, nil
		case "ReadWriteMany", "ReadOnlyMany":
			loc := provider.metadata.location
			parameter[location] = loc
			parameter[storageAccount] = "banzaicloudtest"
			parameter[skuName] = "Standard_LRS"
			return parameter, nil
		}
	}
	return nil, errors.New("could not determine parameters")
}

func determineProvisioner(pvc *v1.PersistentVolumeClaim) (string, error) {
	for _, mode := range pvc.Spec.AccessModes {
		switch mode {
		case "ReadWriteOnce":
			return "kubernetes.io/azure-disk", nil
		case "ReadWriteMany":
			return "kubernetes.io/azure-file", nil
		case "ReadOnlyMany":
			return "kubernetes.io/azure-file", nil
		}
	}
	return "", nil
}