package observer

import "github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"

type State struct {
	ClusterStats *zk.ClusterStats
}
