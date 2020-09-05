package provision

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/sts"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (p *Provision) ProvisionStatefulset() error {
	actual := &appsv1.StatefulSet{}
	name := p.Workload.Name
	namespace := p.Workload.Namespace

	s := sts.NewSTS(p.Workload, p.Labels)

	expect, err := s.GenerateStatefulset()
	if err != nil {
		return err
	}
	if err := controllerutil.SetControllerReference(p.Workload, expect, p.Scheme); err != nil {
		return err
	}

	if err := p.Client.Get(types.NamespacedName{Name: name, Namespace: namespace}, actual); err != nil && errors.IsNotFound(err) {
		p.Log.Info("Creating StatefulSet.",
			"namespace", namespace, "name", name)

		if err := p.Client.Create(expect); err != nil {
			return err
		}

		p.ExpectSts = expect

		msg := fmt.Sprintf(model.MessageZooKeeperStatefulset, name)
		p.Recorder.Event(p.Workload, corev1.EventTypeNormal, model.ZooKeeperStatefulset, msg)

		p.Log.Info("StatefulSet create complete.",
			"namespace", namespace, "name", name)

	} else if err != nil {
		return err
	} else {
		p.ExpectSts = expect
		p.ActualSts = actual
	}

	return nil
}
