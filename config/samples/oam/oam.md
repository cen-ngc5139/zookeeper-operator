# zookeeper-operator access OAM framework

## Prerequisite

  Make sure [`OAM runtime`](https://github.com/crossplane/oam-kubernetes-runtime/blob/master/README.md) was installed and started.

## Install Zookeeper-operator

   
  Step 1. Modify the Makefile file and replace `CRD_OPTIONS ?= "crd:crdVersions=v1,trivialVersions=false"` with CRD_OPTIONS ?= "crd:trivialVersions=true" to generate apiextensions.k8s.io/v1 version CRD.
  
  Step 2. Use `make install` to install zookeeper-operator CRD.
  
  Step 3. Use `make deploy` to install zookeeper-operator Controller.
  
## Registry Zookeeper-operator to OAM Workload

```yaml
apiVersion: core.oam.dev/v1alpha2
kind: WorkloadDefinition
metadata:
  name: workloads.cache.ghostbaby.io
spec:
  definitionRef:
    name: workloads.cache.ghostbaby.io
```

## Create Zookeeper-operator Component
```yaml
apiVersion: core.oam.dev/v1alpha2
kind: Component
metadata:
  name: zk-component
spec:
  workload:
    apiVersion: cache.ghostbaby.io/v1alpha1
    kind: Workload
    spec:
      version: v3.5.6
      cluster:
        name: test
        resources:
          requests:
            cpu: 100m
            memory: 500Mi
        exporter:
          exporter: true
          exporterImage: ghostbaby/zookeeper_exporter
          exporterVersion: v3.5.6
          disableExporterProbes: false
  parameters:
    - name: name
      fieldPaths:
        - metadata.name
```

## Create Zookeeper Cluster and bind ManualScalerTrait
```yaml
apiVersion: core.oam.dev/v1alpha2
kind: ApplicationConfiguration
metadata:
  name: zk-appconfig
spec:
  components:
    - componentName: zk-component
      parameterValues:
        - name: name
          value: ghostbaby
      traits:
        - trait:
            apiVersion: core.oam.dev/v1alpha2
            kind: ManualScalerTrait
            metadata:
              name: zk-appconfig-trait
            spec:
              replicaCount: 3
```