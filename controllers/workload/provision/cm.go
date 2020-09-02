package provision

import (
	"fmt"
	"reflect"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/cm"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (p *Provision) ProvisionConfigMap() error {
	namespace := p.Workload.Namespace

	cm := cm.NewCM(p.Workload, p.Labels, p.ExpectSts, p.ActualSts)

	scs, err := cm.GenerateConfigMap()
	if err != nil {
		return err
	}

	for _, sc := range scs {
		cm := &corev1.ConfigMap{}
		if err := controllerutil.SetControllerReference(p.Workload, sc, p.Scheme); err != nil {
			p.Log.Error(err, "Set OwnerReference fail.", "namespace", cm.Namespace, "name", cm.Name)
			return err
		}

		err = p.Client.Get(types.NamespacedName{Name: sc.Name, Namespace: namespace}, cm)
		if err != nil && errors.IsNotFound(err) {
			p.Log.Info("Creating ZooKeeper ConfigMap .")
			err = p.Client.Create(sc)
			if err != nil {
				p.Log.Error(err, "Create ConfigMap fail.", "namespace", cm.Namespace, "name", cm.Name)
				return err
			}
			msg := fmt.Sprintf(model.MessageZooKeeperConfigMap, sc.Name)
			p.Recorder.Event(p.Workload, corev1.EventTypeNormal, model.ZooKeeperConfigMap, msg)
		} else if err != nil {
			p.Log.Error(err, "Create ConfigMap fail.", "namespace", cm.Namespace, "name", cm.Name)
			return err
		} else {
			if !reflect.DeepEqual(sc.Data, cm.Data) {
				cm.Data = sc.Data
				msg := fmt.Sprintf(model.UpdateMessageZooKeeperConfigMap, sc.Name)
				p.Recorder.Event(p.Workload, corev1.EventTypeNormal, model.ZooKeeperConfigMap, msg)
				p.Log.Info("Updating ConfigMap .", "namespace", cm.Namespace, "name", cm.Name)
				err = p.Client.Update(cm)
				if err != nil {
					p.Log.Error(err, "Update ConfigMap fail.", "namespace", cm.Namespace, "name", cm.Name)
					return err
				}
			}
		}

	}

	return nil
}
