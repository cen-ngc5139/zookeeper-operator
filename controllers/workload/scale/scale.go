package scale

import (
	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Scale struct {
	Workload  *cachev1alpha1.Workload
	Client    k8s.Client
	Recorder  record.EventRecorder
	Log       logr.Logger
	Labels    map[string]string
	Scheme    *runtime.Scheme
	ExpectSts *appsv1.StatefulSet
	ActualSts *appsv1.StatefulSet
}

func (s *Scale) Reconcile() error {

	if err := s.StatefulSet(); err != nil {
		return err
	}

	if err := s.UpScale(); err != nil {
		return err
	}

	if err := s.ReConfig(); err != nil {
		return err
	}

	if err := s.DownScale(); err != nil {
		return err
	}

	return nil
}
