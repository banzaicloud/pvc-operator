# pvc-operator

### Introduction

This operator helps to use [Kubernetes Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) easier on some cloud providers, by dynamically
creating the required accounts, classes and more. It allows to use exactly the same [helm](https://helm.sh) chart for all the
supported providers.

Currently supported Providers/StorageClasses:

- Azure(AKS)
    - AzureFile
    - AzureDisk
    - Nfs
    
- Amazon
    - AWSElasticBlockStore
    - Nfs
    
- Google
    - GCEPersistentDisk
    - Nfs
    
### Installation

The user should use the [operator.yaml](https://github.com/banzaicloud/pvc-operator/blob/master/deploy/operator.yaml)
to deploy the operator to a Kubernetes cluster. If the cluster uses [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/)
deploy the [rbac.yaml](https://github.com/banzaicloud/pvc-operator/blob/master/deploy/rbac.yaml) first and add the following line after
the specs in the `operator.yaml`.

```serviceAccountName: pvc-operator```

### Cloud Specific Requirements

In case of `AzureFile` a Storage Account needs to be created. The operator handles this creation automatically
but some permissions has to be set.

- MSI has to be [enabled](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm#enable-msi-on-your-vm)
- Grant Access to your VMs to [create](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm#grant-your-vm-access-to-a-resource-group-in-azure-resource-manager) a Storage Account
Instead of adding `Read` role use the `Storage Account Owner`.

### Usage

The given chart should include a `Persistent Volume Claim` which includes a [StorageClass](https://kubernetes.io/docs/concepts/storage/storage-classes/)
name and an `Access Mode`. If the chosen Access Mode is supported on the required cloud provider
the operator will create a proper `StorageClass`. This class will be reused by other charts as well.

### FAQ

#### 1. How does this project uses Kubernetes Namespaces?

Right now every instance created by this `Operator` is in `default` namespace, but we are planning to add support
for different namespaces.

#### 2. How is the cloud type determined?

To determine cloud type we use the metadata server accessible from every virtual machine inside the cloud.  

#### 3. Do I need to add my cloud related credentials to this project?

No, this project determines every credentials from the `metadata` server, and uses the assumption that you have
a couple of right because this `Operator` runs inside a Kubernetes cluster.


#### 4. How come that the Nfs type StorageClass can be used in any cloud provider?

The Nfs StorageClass solution is based on the [Kubernetes external-storages](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
incubator github project. It uses an own Nfs server deployment which requires a `ReadWriteOnly` Kubernetes Persistent Volume
and shares this across the cluster.

#### 5. Why this project exists?

Here at [banzaicloud](www.banzaicloud.com) we try to automate everything so you don't have to. If you need `Kubernetes`
cluster please check out [Pipeline](github.com/banzaicloud/pipeline). You already have one, but you are struggling
to find and configure the right `Persistent Storage` for your needs? Then this project is for you. 

#### 6. What's next for this project for the near future?

- The priority will be to support more `ReadWriteMany` StorageClass on all providers. We prefer to build on cloud
specific storage solutions, but more generic solutions will come as well.

- We would like to add support to create a object storage inside a Kubernetes Cluster. This approach will help in case of
[Spark History Server](https://banzaicloud.com/blog/spark-history-server-cloud/)

#### 7. Is this project production ready?

No, not yet. We need to introduce at least some unit tests and maybe an integration test too.