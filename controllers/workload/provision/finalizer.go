package provision

import (
	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/finalizer"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"
)

func (p *Provision) FinalizersFor(
	zk *cachev1alpha1.Workload,
) []finalizer.Finalizer {
	clusterName := utils.ExtractNamespacedName(zk)
	return []finalizer.Finalizer{
		p.Observers.Finalizer(clusterName),
	}
}
