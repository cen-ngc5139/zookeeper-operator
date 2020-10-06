package workload

import (
	"context"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/prometheus"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/finalizer"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/observer"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"

	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// ReconcileWorkload implement the Reconciler interface and lcm.Controller interface.
type ReconcileWorkload struct {
	Getter
	Workload      *cachev1alpha1.Workload
	CTX           context.Context
	Client        k8s.Client
	Recorder      record.EventRecorder
	Log           logr.Logger
	DClient       k8s.DClient
	Scheme        *runtime.Scheme
	Observers     *observer.Manager
	Monitor       *prometheus.GenericClientset
	Labels        map[string]string
	ZKClient      *zk.BaseClient
	ObservedState *observer.State
	Finalizers    finalizer.Handler
}

func (w *ReconcileWorkload) Reconcile() error {
	w.Client.WithContext(w.CTX)
	option := w.GetOptions()

	if err := w.ProvisionWorkload(w.CTX, w.Workload, option).Reconcile(); err != nil {
		return err
	}

	if !w.Workload.GetDeletionTimestamp().IsZero() {
		return nil
	}

	if err := w.ScaleWorkload(w.CTX, w.Workload, option).Reconcile(); err != nil {
		return err
	}

	if err := w.RolloutWorkload(w.CTX, w.Workload, option).Reconcile(); err != nil {
		return err
	}
	return nil
}
