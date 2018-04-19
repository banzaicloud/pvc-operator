package providers

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/coreos/operator-sdk/pkg/sdk/query"
)

type CommonProvider interface {
	CreateStorageClass(*v1.PersistentVolumeClaim) error
	GenerateMetadata() error
}

func DetermineProvider() (CommonProvider, error) {
	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2017-12-01", nil)
	if err != nil {
		logrus.Errorf("Could not create a proper http request %s", err.Error())
		return nil, err
	}
	req.Header.Set("Metadata", "true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &AzureProvider{}, nil
}

func CheckStorageClassExistence(name string) bool {
	storageClassList := &storagev1.StorageClassList{
		TypeMeta: v12.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
	}
	if err := query.List("default", storageClassList); err != nil {
		logrus.Infof("Error happened during listing storageclass %s", err.Error())
		return false
	}
	for _, storageClass := range storageClassList.Items {
		if storageClass.Name == name {
			logrus.Info("Storageclass exist!")
			return true
		}
	}
	logrus.Info("Storageclass does not exist!")
	return false
}
