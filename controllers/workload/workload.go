package workload

import (
	"context"

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
	ExpectCR      *cachev1alpha1.Workload
	ActualCR      *cachev1alpha1.Workload
	Labels        map[string]string
	ZKClient      *zk.BaseClient
	ObservedState *observer.State
	Finalizers    finalizer.Handler
}

func (w *ReconcileWorkload) Reconcile() error {
	w.Client.WithContext(w.CTX)
	option := w.GetOptions()

	if err := w.Provision(option); err != nil {
		return err
	}
	return nil
}

func (w *ReconcileWorkload) Provision(option *GetOptions) error {
	return w.ProvisionWorkload(w.CTX, w.Workload, option).Reconcile()
}

func (w *ReconcileWorkload) Delete() error {
	return nil
}

func (w *ReconcileWorkload) Scale() error {
	return nil
}

func (w *ReconcileWorkload) Update() error {
	return nil
}
