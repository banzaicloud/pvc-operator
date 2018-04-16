// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"bytes"
	"testing"
)

const (
	// test constants for app-operator
	appRepoPath   = "github.com/example-inc/app-operator"
	appKind       = "App"
	appApiDirName = "app"
	appAPIVersion = appGroupName + "/" + appVersion
	appVersion    = "v1alpha1"
	appGroupName  = "app.example.com"
)

const mainExp = `package main

import (
	"context"
	"runtime"

	stub "github.com/example-inc/app-operator/pkg/stub"
	sdk "github.com/coreos/operator-sdk/pkg/sdk"
	sdkVersion "github.com/coreos/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()
	sdk.Watch("app.example.com/v1alpha1", "App", "default", 5)
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}
`

func TestGenMain(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderMainFile(buf, appRepoPath, appAPIVersion, appKind); err != nil {
		t.Error(err)
		return
	}

	if mainExp != buf.String() {
		t.Errorf("want %v, got %v", mainExp, buf.String())
	}
}

const handlerExp = `package stub

import (
	"github.com/example-inc/app-operator/pkg/apis/app/v1alpha1"

	"github.com/coreos/operator-sdk/pkg/sdk/action"
	"github.com/coreos/operator-sdk/pkg/sdk/handler"
	"github.com/coreos/operator-sdk/pkg/sdk/types"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() handler.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx types.Context, event types.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.App:
		err := action.Create(newbusyBoxPod(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create busybox pod : %v", err)
			return err
		}
	}
	return nil
}

// newbusyBoxPod demonstrates how to create a busybox pod
func newbusyBoxPod(cr *v1alpha1.App) *v1.Pod {
	labels := map[string]string{
		"app": "busy-box",
	}
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "busy-box",
			Namespace:    "default",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "App",
				}),
			},
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
`

func TestGenHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderHandlerFile(buf, appRepoPath, appKind, appApiDirName, appVersion); err != nil {
		t.Error(err)
		return
	}
	if handlerExp != buf.String() {
		t.Errorf("want %v, got %v", handlerExp, buf.String())
	}
}

const registerExp = `package v1alpha1

import (
	sdkK8sutil "github.com/coreos/operator-sdk/pkg/util/k8sutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	version   = "v1alpha1"
	groupName = "play.example.com"
)

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
	// SchemeGroupVersion is the group version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: groupName, Version: version}
)

func init() {
	sdkK8sutil.AddToSDKScheme(AddToScheme)
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&PlayService{},
		&PlayServiceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
`

func TestGenRegister(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderAPIRegisterFile(buf, "PlayService", "play.example.com", "v1alpha1"); err != nil {
		t.Error(err)
		return
	}
	if registerExp != buf.String() {
		t.Errorf("want %v, got %v", registerExp, buf.String())
	}
	// TODO add verification
}

const typesExp = `package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PlayServiceList struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`\n" +
	`	metav1.ListMeta ` + "`" + `json:"metadata"` + "`\n" +
	`	Items           []PlayService ` + "`" + `json:"items"` + "`" + `
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PlayService struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`\n" +
	`	metav1.ObjectMeta ` + "`" + `json:"metadata"` + "`\n" +
	`	Spec              PlayServiceSpec   ` + "`" + `json:"spec"` + "`\n" +
	`	Status            PlayServiceStatus ` + "`" + `json:"status,omitempty"` + "`" + `
}

type PlayServiceSpec struct {
	// Fills me
}
type PlayServiceStatus struct {
	// Fills me
}
`

func TestGenTypes(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderAPITypesFile(buf, "PlayService", "v1alpha1"); err != nil {
		t.Error(err)
		return
	}
	if typesExp != buf.String() {
		t.Errorf("want %v, got %v", typesExp, buf.String())
	}
}

const buildExp = `#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go > /dev/null; then
	echo "golang needs to be installed"
	exit 1
fi

BIN_DIR="$(pwd)/tmp/_output/bin"
mkdir -p ${BIN_DIR}
PROJECT_NAME="play"
REPO_PATH="github.com/coreos/play"
BUILD_PATH="${REPO_PATH}/cmd/${PROJECT_NAME}"
echo "building "${PROJECT_NAME}"..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BIN_DIR}/${PROJECT_NAME} $BUILD_PATH
`

const dockerFileExp = `FROM alpine:3.6

ADD tmp/_output/bin/play /usr/local/bin/play

RUN adduser -D play
USER play
`

func TestGenBuild(t *testing.T) {
	buf := &bytes.Buffer{}
	projectName := "play"
	if err := renderBuildFile(buf, "github.com/coreos/play", projectName); err != nil {
		t.Error(err)
		return
	}
	if buildExp != buf.String() {
		t.Errorf("want %v, got %v", buildExp, buf.String())
	}

	buf = &bytes.Buffer{}
	if err := renderDockerBuildFile(buf); err != nil {
		t.Error(err)
		return
	}
	if dockerBuildTmpl != buf.String() {
		t.Errorf("want %v, got %v", dockerBuildTmpl, buf.String())
	}

	buf = &bytes.Buffer{}
	if err := renderDockerFile(buf, projectName); err != nil {
		t.Error(err)
		return
	}
	if dockerFileExp != buf.String() {
		t.Errorf("want %v, got %v", dockerFileExp, buf.String())
	}
}

const boilerplateExp = `
`

const updateGeneratedExp = `#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DOCKER_REPO_ROOT="/go/src/github.com/coreos/play"
IMAGE=${IMAGE:-"gcr.io/coreos-k8s-scale-testing/codegen"}

docker run --rm \
  -v "$PWD":"$DOCKER_REPO_ROOT" \
  -w "$DOCKER_REPO_ROOT" \
  "$IMAGE" \
  "/go/src/k8s.io/code-generator/generate-groups.sh"  \
  "deepcopy" \
  "github.com/coreos/play/pkg/generated" \
  "github.com/coreos/play/pkg/apis" \
  "play:v1alpha1" \
  --go-header-file "./tmp/codegen/boilerplate.go.txt" \
  $@
`

func TestCodeGen(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderBoilerplateFile(buf, "play"); err != nil {
		t.Error(err)
		return
	}
	if boilerplateExp != buf.String() {
		t.Errorf("want %v, got %v", boilerplateExp, buf.String())
	}

	buf = &bytes.Buffer{}
	if err := renderUpdateGeneratedFile(buf, "github.com/coreos/play", "play", "v1alpha1"); err != nil {
		t.Error(err)
		return
	}
	if updateGeneratedExp != buf.String() {
		t.Errorf("want %v, got %v", updateGeneratedExp, buf.String())
	}
}

func TestGenGopkg(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderGopkgTomlFile(buf); err != nil {
		t.Error(err)
		return
	}

	if gopkgTomlTmpl != buf.String() {
		t.Errorf("want %v, got %v", gopkgTomlTmpl, buf.String())
	}

	buf = &bytes.Buffer{}
	if err := renderGopkgLockFile(buf); err != nil {
		t.Error(err)
		return
	}
	if gopkgLockTmpl != buf.String() {
		t.Errorf("want %v, got %v", gopkgLockTmpl, buf.String())
	}
}

const configExp = `apiVersion: app.example.com/v1alpha1
kind: AppService
projectName: app-operator
`

func TestGenConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := renderConfigFile(buf, "app.example.com/v1alpha1", "AppService", "app-operator"); err != nil {
		t.Error(err)
	}
	if configExp != buf.String() {
		t.Errorf("want %v, got %v", configExp, buf.String())
	}
}

const operatorYamlExp = `apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: appservices.app.example.com
spec:
  group: app.example.com
  names:
    kind: AppService
    listKind: AppServiceList
    plural: appservices
    singular: appservice
  scope: Namespaced
  version: v1alpha1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: app-operator
  template:
    metadata:
      labels:
        name: app-operator
    spec:
      containers:
        - name: app-operator
          image: quay.io/coreos/operator-sdk-dev:app-operator
          command:
          - app-operator
          imagePullPolicy: Always
`

const rbacYamlExp = `kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: app-operator
rules:
- apiGroups:
  - "*"
  resources:
  - "*"
  verbs:
  - "*"

---

kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: default-account-app-operator
subjects:
- kind: ServiceAccount
  name: default
roleRef:
  kind: Role
  name: app-operator
  apiGroup: rbac.authorization.k8s.io
`

func TestGenDeploy(t *testing.T) {
	buf := &bytes.Buffer{}
	projectName := "app-operator"
	if err := renderOperatorYaml(buf, "AppService", "app.example.com/v1alpha1", projectName, "quay.io/coreos/operator-sdk-dev:app-operator"); err != nil {
		t.Error(err)
	}
	if operatorYamlExp != buf.String() {
		t.Errorf("want %v, got %v", operatorYamlExp, buf.String())
	}

	buf = &bytes.Buffer{}
	if err := renderRBACYaml(buf, projectName); err != nil {
		t.Error(err)
	}
	if rbacYamlExp != buf.String() {
		t.Errorf("want %v, got %v", rbacYamlExp, buf.String())
	}

}
