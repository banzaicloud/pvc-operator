# pvc-operator

### Introduction

This operator helps to use [Kubernetes Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) easier on some cloud providers, by dynamically
creating the required accounts, classes and more. It allows to use exactly the same [helm](https://helm.sh) chart for all the
supported providers.

Currently supported Providers/StorageClasses:

- Azure(AKS)
    - AzureFile
    - AzureDisk
    
- Amazon
    - AWSElasticBlockStore
    
- Google
    - GCEPersistentDisk
    
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

### Future Work

- Add support for more `ReadWriteMany` StorageClass eg.: Glusterfs, CephFS
- Add support to create a Blob storage from a Kubernetes Cluster for example for [Spark History Server](https://banzaicloud.com/blog/spark-history-server-cloud/).