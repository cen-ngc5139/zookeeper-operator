## ZooKeeper Operator
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://godoc.org/github.com/apache/rocketmq-operator/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ghostbaby/zookeeper-operator)](https://goreportcard.com/report/github.com/Ghostbaby/zookeeper-operator)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/Ghostbaby/zookeeper-operator.svg)](http://isitmaintained.com/project/Ghostbaby/zookeeper-operator "Average time to resolve an issue")
[![Percentage of issues still open](http://isitmaintained.com/badge/open/Ghostbaby/zookeeper-operator.svg)](http://isitmaintained.com/project/Ghostbaby/zookeeper-operator "Percentage of issues still open")
    
## Overview
ZooKeeper Operator is to manage ZooKeeper service instances deployed on the Kubernetes cluster.
It is built using the [Kubebuilder SDK](https://kubebuilder.io/), which is part of the [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder).

## Features

With this operator, you're able to deploy and manage a HA Zookeeper Cluster:

- [x] Provision a Zookeeper cluster in a scalable and high-available way.
- [x] Update the spec of the deployed Zookeeper cluster to do adjustments like replicas (scalability) and resources.
  - ScaleUp
  - ScaleDown
  - Rollout
  - Observe
- [x] Create Prometheus target for the Zookeeper node.
- [x] Delete the Zookeeper cluster and all the related resources owned by the instance.

## Design

Diagram below shows the overall design of this operator,

![架构图](https://raw.githubusercontent.com/Ghostbaby/picgo/master/image2019-11-21_15-38-17.png)


For more design details, check the [architecture](design/zookeeper-operator-cn.md) document.

## Installation

You can follow the [installation guide](docs/installation.md) to deploy this operator to your K8s clusters.

Additionally, follow [sample deployment guide](./docs/sample_deploy_guide.md) to have a try of deploying the sample to your K8s clusters.

## Versioning & Dependencies

| Component \ Versions |  0.5.0 | 1.0.0 | 1.1.0 |
|----------------------|--------|-------|-------|
| **Zookeeper**        | 3.5.6  | [TBD] | [TBD] |
|                      |        |               |
| agent                | 0.0.1  | [TBD] | [TBD] |

## Compatibilities

| Kubernetes / Versions |  0.5.0  |  1.0.0  | 1.1.0 |
|-----------------------|---------|---------|------|
|     1.17              |    +    | [TBD] | [TBD] |
|     1.18              |    +    | [TBD] | [TBD] |
|     1.19              |    +    | [TBD] | [TBD] |

**Notes:** `+`= verified `-`= not verified

## Development

Interested in contributions? Follow the [CONTRIBUTING](./docs/DEVELOPMENT.md) guide to start on this project. Your contributions will be highly appreciated and creditable.

## Community

* Send mail to Ghostbaby mail:  zhuhuijunzhj@gmail.com

## Documents

See documents [here](./docs).

## Additional Documents

* [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
* [Kubebuilder](https://book.kubebuilder.io/)
* [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)


## License

[Apache-2.0](https://github.com/goharbor/harbor-cluster-operator/blob/master/LICENSE)
