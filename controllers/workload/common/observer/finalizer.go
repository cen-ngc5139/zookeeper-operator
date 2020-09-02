package observer

import (
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/finalizer"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// FinalizerName registered for each elasticsearch resource
	FinalizerName = "finalizer.zookeeper.ymmoa.inc/observer"
)

// Finalizer returns a finalizer to be executed upon deletion of the given cluster,
// that makes sure the cluster is not observed anymore
func (m *Manager) Finalizer(cluster types.NamespacedName) finalizer.Finalizer {
	return finalizer.Finalizer{
		Name: FinalizerName,
		Execute: func() error {
			m.StopObserving(cluster)
			return nil
		},
	}
}
