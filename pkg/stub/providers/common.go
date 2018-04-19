package providers

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/coreos/operator-sdk/pkg/sdk/query"
	"fmt"
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
		"aws": "http://169.254.169.254/latest/meta-data/",
	}
	for key, value := range providers{
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
		if resp.StatusCode == 404 {
			continue
		}
		switch key {
		case "azure":
			return &AzureProvider{}, nil
		case "aws":
			return &AwsProvider{}, nil
		}

	}
	return nil, fmt.Errorf("could not determine cloud provider")
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
