package main

import (
	"context"
	"runtime"

	"github.com/banzaicloud/pvc-operator/pkg/stub"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()
	resyncPeriod := 0
	sdk.Watch("banzaicloud.com/v1alpha1", "ObjectStore", metav1.NamespaceAll, resyncPeriod)
	sdk.Watch("v1", "PersistentVolumeClaim", metav1.NamespaceAll, resyncPeriod)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}
