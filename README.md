## PVC Operator

This operator builds upon the new [Kubernetes operators SDK](https://banzaicloud.com/blog/operator-sdk/) and used on different cloud providers.

### Introduction

This operator makes using [Kubernetes Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) easier on cloud providers, by dynamically creating the required accounts, classes and more. It allows to use exactly the same [Helm](https://helm.sh) chart on all the supported providers, there is no need to create cloud specific Helm charts.

The Currently supported list of Providers/StorageClasses are:

- Azure(AKS)
    - AzureFile
    - AzureDisk
    - NFS
    
- Amazon
    - AWSElasticBlockStore
    - NFS
    
- Google
    - GCEPersistentDisk
    - NFS
    
### Installation

- Branch `0.0.2` contains an earlier version of this operator if you want to use this please follow this guide:
The user should use the [operator.yaml](https://github.com/banzaicloud/pvc-operator/blob/0.0.2/deploy/operator.yaml)
to deploy the operator to a Kubernetes cluster. If the cluster uses [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) deploy the [rbac.yaml](https://github.com/banzaicloud/pvc-operator/blob/0.0.2/deploy/rbac.yaml) first and add the following line after the specs in the `operator.yaml`.

```serviceAccountName: pvc-operator```

- In case of `master` branch, please use the [crd.yaml](https://github.com/banzaicloud/pvc-operator/blob/master/deploy/crd.yaml) first then
deploy the operator itself by using the [operator.yaml](https://github.com/banzaicloud/pvc-operator/blob/master/deploy/operator.yaml).
If the cluster uses [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/) deploy the [rbac.yaml](https://github.com/banzaicloud/pvc-operator/blob/master/deploy/rbac.yaml).

### Cloud Specific Requirements

In case of `AzureFile` a Storage Account needs to be created. The operator handles the creation automatically
but some permissions have to be set.

- MSI has to be [enabled](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm#enable-msi-on-your-vm)
- Grant Access to your VMs to [create](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm#grant-your-vm-access-to-a-resource-group-in-azure-resource-manager) a Storage Account
Instead of adding `Read` role use the `Storage Account Owner`.

### Usage

The given chart should include a `Persistent Volume Claim` which includes a [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/) name and an `Access Mode`. If the chosen Access Mode is supported on the required cloud provider the operator will create a proper `StorageClass`. This class will be reused by other charts as well.

### FAQ

#### 1. How does this project uses Kubernetes Namespaces?

The operator supports multiple Kubernetes Namespaces. In case of NFS the Deployment will still be created in the `default` namespace. 

#### 2. How is the cloud provider determined?

To determine the cloud provider we use the `metadata` server accessible from every virtual machine within the cloud.  

#### 3. Do I need to add my cloud related credentials to this project?

No need, this project determines every credentials from the `metadata` server and uses the assumption that you have
a couple of right because this `Operator` runs inside a Kubernetes cluster (IAM or instance profile roles, etc).

#### 4. How come that the NFS type StorageClass can be used in any cloud provider?

The NFS StorageClass solution is based on the [Kubernetes external-storages](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs) project. It uses it's own NFS server deployment which requires a `ReadWriteOnly` Kubernetes Persistent Volume and shares this across the cluster.

#### 5. Why this project exists?

Here at [Banzai Cloud](https://banzaicloud.com) we try to automate everything so you don't have to. If you need to provision a Kubernetes cluster please check out [Pipeline](github.com/banzaicloud/pipeline). If you already have one, but you are struggling to find and configure the right `Persistent Storage` for your need then this project is for you. 

#### 6. What's next for this project for the near future?

- The priority will be to support more `ReadWriteMany` StorageClasses on all providers. We prefer to build on cloud
specific storage solutions, but more generic solutions will come as well.

- We would like to add support to create an object storage inside a Kubernetes Cluster. This approach will help e.g in case of
[Spark History Server](https://banzaicloud.com/blog/spark-history-server-cloud/).

#### 7. Is this project production ready?

Although we use it internally, it's not yet. We need to add at least some unit tests and potentially an integration test too.
