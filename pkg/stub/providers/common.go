package providers

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
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
