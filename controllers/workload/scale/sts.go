package scale

import (
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/sts"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (s *Scale) StatefulSet() error {
	actual := &appsv1.StatefulSet{}
	name := s.Workload.Name
	namespace := s.Workload.Namespace

	newSts := sts.NewSTS(s.Workload, s.Labels)

	expect, err := newSts.GenerateStatefulset()
	if err != nil {
		return err
	}
	if err := controllerutil.SetControllerReference(s.Workload, expect, s.Scheme); err != nil {
		return err
	}

	if err := s.Client.Get(types.NamespacedName{Name: name, Namespace: namespace}, actual); err != nil && errors.IsNotFound(err) {
		s.ExpectSts = expect
	} else if err != nil {
		return err
	} else {
		s.ExpectSts = expect
		s.ActualSts = actual
	}

	return nil
}
