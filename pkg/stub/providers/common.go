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
		"aws":   "http://169.254.169.254/latest/meta-data/",
		"google": "http://169.254.169.254/0.1/meta-data/",
	}
	for key, value := range providers {
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
		if resp.StatusCode == 404 || resp.StatusCode == 405 {
			continue
		}
		switch key {
		case "azure":
			return &AzureProvider{}, nil
		case "aws":
			return &AwsProvider{}, nil
		case "google":
			return &GoogleProvider{}, nil
		}

	}
	return nil, fmt.Errorf("could not determine cloud provider")
}

func CheckStorageClassExistence(name string) bool {
	storageClass := &storagev1.StorageClass{
		TypeMeta: v12.TypeMeta{
			Kind:       "StorageClass",
			APIVersion: "storage.k8s.io/v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name: name,
		},
	}
	if err := query.Get(storageClass); err != nil {
		logrus.Infof("Storageclass does not exist %s", err.Error())
		return false
	}
	logrus.Infof("Storageclass %s exists!", name)
	return true
}
