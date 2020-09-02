package controllers

import (
	"context"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	w "github.com/ghostbaby/zookeeper-operator/controllers/workload"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Reconciler interface {
	// Reconcile the dependent service.
	Reconcile() error
}

type ServiceGetter interface {
	// For Workload
	Workload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler
}

type GetOptions struct {
	Client   k8s.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	DClient  k8s.DClient
	Scheme   *runtime.Scheme
	Labels   map[string]string
}

type ServiceGetterImpl struct {
}

func (impl *ServiceGetterImpl) Workload(ctx context.Context, workload *cachev1alpha1.Workload, options *GetOptions) Reconciler {
	return &w.ReconcileWorkload{
		Workload: workload,
		Client:   options.Client,
		Recorder: options.Recorder,
		Log:      options.Log,
		DClient:  options.DClient,
		Scheme:   options.Scheme,
		CTX:      ctx,
		Labels:   options.Labels,
		Getter:   &w.GetterImpl{},
	}
}
