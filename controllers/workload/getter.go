package workload

import (
	"context"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/prometheus"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/rollout"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/scale"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/finalizer"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/observer"

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
	ProvisionWorkload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler

	// For Scale
	ScaleWorkload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler

	// For Rollout
	RolloutWorkload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler
}

type GetOptions struct {
	Client        k8s.Client
	Recorder      record.EventRecorder
	Log           logr.Logger
	DClient       k8s.DClient
	Scheme        *runtime.Scheme
	Labels        map[string]string
	Observers     *observer.Manager
	ZKClient      *zk.BaseClient
	ObservedState *observer.State
	Monitor       *prometheus.GenericClientset
	Finalizers    finalizer.Handler
	ExpectSts     *appsv1.StatefulSet
	ActualSts     *appsv1.StatefulSet
}

type GetterImpl struct {
}

func (w *ReconcileWorkload) GetOptions() *GetOptions {
	return &GetOptions{
		Client:        w.Client,
		Recorder:      w.Recorder,
		Log:           w.Log,
		DClient:       w.DClient,
		Scheme:        w.Scheme,
		Labels:        w.Labels,
		Observers:     w.Observers,
		ZKClient:      w.ZKClient,
		ObservedState: w.ObservedState,
		Finalizers:    w.Finalizers,
		Monitor:       w.Monitor,
	}
}

func (impl *GetterImpl) ProvisionWorkload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler {
	return &provision.Provision{
		Workload:      workload,
		CTX:           ctx,
		Client:        options.Client,
		Recorder:      options.Recorder,
		Log:           options.Log,
		Labels:        options.Labels,
		Scheme:        options.Scheme,
		Observers:     options.Observers,
		ZKClient:      options.ZKClient,
		ObservedState: options.ObservedState,
		Finalizers:    options.Finalizers,
		Monitor:       options.Monitor,
	}
}

func (impl *GetterImpl) ScaleWorkload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler {
	return &scale.Scale{
		Workload: workload,
		Client:   options.Client,
		Recorder: options.Recorder,
		Log:      options.Log,
		Labels:   options.Labels,
		Scheme:   options.Scheme,
	}
}

func (impl *GetterImpl) RolloutWorkload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler {
	return &rollout.Rollout{
		Workload: workload,
		Client:   options.Client,
		Recorder: options.Recorder,
		Log:      options.Log,
		Labels:   options.Labels,
		Scheme:   options.Scheme,
	}
}
