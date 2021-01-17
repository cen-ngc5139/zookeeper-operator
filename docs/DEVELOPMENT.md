# Development

This doc explains how to set up a development environment, so you can get started
contributing to `zookeeper-operator` or build a PoC (Proof of Concept). 

## Prerequisites

1. Golang version 1.13+
2. Kubernetes version v1.15+ with `~/.kube/config` configured.
4. Kustomize version 3.8+
5. Kubebuilder version 2.0+

## Build
* Clone this project

```shell script
git clone git@github.com:Ghostbaby/zookeeper-operator.git
```

* Install Zookeeper CRD into your cluster

```shell script
make install
```

## Develop & Debug
If you change Zookeeper CRD, remember to rerun `make install`.

Use the following command to develop and debug.

```shell script
$ make run
```

For example, use the following command to create an zookeeper cluster.
```shell script
$ cd ./config/samples

$ kubectl apply -f cache_v1alpha1_workload.yaml
workload.zk.cache.ghostbaby.io/workload-sample created

$ kubectl get pods -n pg
workload-sample-0                   3/3     Running   0          3m
workload-sample-1                   3/3     Running   0          3m
workload-sample-2                   3/3     Running   0          3m
```

## Make a pull request
Remember to write unit-test and e2e test before making a pull request.
