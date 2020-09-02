package workload

import (
	"context"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/provision"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Reconciler interface {
	// Reconcile the dependent service.
	Reconcile() error
}

type Getter interface {
	// For Provision
	ProvisionW(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler
}

type GetOptions struct {
	Client   k8s.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	DClient  k8s.DClient
	Scheme   *runtime.Scheme
	Labels   map[string]string
}

type GetterImpl struct {
}

func (impl *GetterImpl) ProvisionW(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler {
	return &provision.Provision{
		Workload: workload,
		CTX:      ctx,
		Client:   options.Client,
		Recorder: options.Recorder,
		Log:      options.Log,
		Labels:   options.Labels,
		Scheme:   options.Scheme,
	}
}

func (w *ReconcileWorkload) GetOptions() *GetOptions {
	return &GetOptions{
		Client:   w.Client,
		Recorder: w.Recorder,
		Log:      w.Log,
		DClient:  w.DClient,
		Scheme:   w.Scheme,
		Labels:   w.Labels,
	}
}
