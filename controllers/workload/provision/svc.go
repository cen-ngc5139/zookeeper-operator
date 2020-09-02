package provision

import (
	"fmt"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"

	svc2 "github.com/ghostbaby/zookeeper-operator/controllers/workload/common/svc"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// EnsureService makes sure the mongodb statefulset exists
func (p *Provision) ProvisionService() error {
	service := &corev1.Service{}
	var svcs []*corev1.Service
	name := p.Workload.GetName()
	namespace := p.Workload.GetNamespace()

	s := svc2.NewSVC(p.Workload, p.Labels)

	svc := s.GenerateService(name, "ClusterIP")
	svcs = append(svcs, svc)

	svcHeadless := s.GenerateService(name+"-s", "Headless")
	svcs = append(svcs, svcHeadless)

	for _, dep := range svcs {
		if err := controllerutil.SetControllerReference(p.Workload, dep, p.Scheme); err != nil {
			p.Log.Error(err, "SVC set ownerReference fail.", "namespace", dep.Namespace, "name", dep.Name)
			return err
		}
		err := p.Client.Get(types.NamespacedName{Name: dep.Name, Namespace: namespace}, service)
		if err != nil && errors.IsNotFound(err) {
			p.Log.Info("Creating Service.", "namespace", namespace, "name", dep.Name)
			err = p.Client.Create(dep)
			if err != nil {
				return err
			}
			msg := fmt.Sprintf(model.MessageZooKeeperService, name)
			p.Recorder.Event(p.Workload, corev1.EventTypeNormal, model.ZooKeeperService, msg)
			p.Log.Info("Service create complete.", "namespace", namespace, "name", dep.Name)
		}
	}

	return nil
}
