package provision

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/runtime"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
)

type Provision struct {
	Workload  *cachev1alpha1.Workload
	CTX       context.Context
	Client    k8s.Client
	Recorder  record.EventRecorder
	Log       logr.Logger
	Labels    map[string]string
	Scheme    *runtime.Scheme
	ExpectSts *appsv1.StatefulSet
	ActualSts *appsv1.StatefulSet
}

func (p *Provision) Reconcile() error {
	if err := p.ProvisionStatefulset(); err != nil {
		return err
	}

	if err := p.ProvisionConfigMap(); err != nil {
		return err
	}

	if err := p.ProvisionService(); err != nil {
		return err
	}
	return nil
}
