package rollout

import (
	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	commonsts "github.com/ghostbaby/zookeeper-operator/controllers/workload/common/sts"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Rollout struct {
	Workload  *cachev1alpha1.Workload
	Client    k8s.Client
	Recorder  record.EventRecorder
	Log       logr.Logger
	Labels    map[string]string
	Scheme    *runtime.Scheme
	ExpectSts *appsv1.StatefulSet
	ActualSts *appsv1.StatefulSet
}

func (r *Rollout) Reconcile() error {

	expectSts, actualSts, err := commonsts.GetStatefulSet(r.Client, r.Workload, r.Labels, r.Scheme)
	if err != nil {
		return err
	}

	r.ExpectSts = expectSts
	r.ActualSts = actualSts

	if err := r.RollingUpgrades(); err != nil {
		return err
	}

	return nil
}
