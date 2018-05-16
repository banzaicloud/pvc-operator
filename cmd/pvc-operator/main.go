package main

import (
	"context"
	"runtime"

	"github.com/banzaicloud/pvc-operator/pkg/stub"
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
	//namespace, err := k8sutil.GetWatchNamespace()
	//if err != nil {
	//	logrus.Fatalf("Failed to get watch namespace: %v", err)
	//}
	resyncPeriod := 0
	sdk.Watch("banzaicloud.com/v1alpha1", "ObjectStore", "default", resyncPeriod)
	sdk.Watch("v1", "PersistentVolumeClaim", "default", resyncPeriod)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}
