package provision

import (
	"fmt"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/monitor"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (p *Provision) ProvisionMonitor() error {

	name := p.Workload.Name
	namespace := p.Workload.Namespace

	m := monitor.NewMonitor(p.Workload, p.Labels)
	sm, err := m.GenerateMongodbServiceMonitor()
	if err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(p.Workload, sm, p.Scheme); err != nil {
		return err
	}

	if _, err := p.Monitor.PoClient.MonitoringV1().ServiceMonitors(namespace).Get(p.CTX, name, metav1.GetOptions{}); err != nil {
		p.Log.Info("Creating ServiceMonitor %s/%s\n", p.Workload.Namespace, p.Workload.Name)

		_, err := p.Monitor.PoClient.MonitoringV1().ServiceMonitors(namespace).Create(p.CTX, sm, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		msg := fmt.Sprintf(model.MessageZooKeeperServiceMonitor, name)
		p.Recorder.Event(p.Workload, corev1.EventTypeNormal, model.ZooKeeperServiceMonitor, msg)
	}

	return nil
}
