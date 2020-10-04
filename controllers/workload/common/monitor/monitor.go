package monitor

import (
	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	monitorV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Monitor struct {
	Workload *cachev1alpha1.Workload
	Labels   map[string]string
}

func NewMonitor(workload *cachev1alpha1.Workload, labels map[string]string) *Monitor {
	return &Monitor{
		Workload: workload,
		Labels:   labels,
	}
}

func (m *Monitor) GenerateMongodbServiceMonitor() (*monitorV1.ServiceMonitor, error) {
	name := m.Workload.Name
	namespace := m.Workload.Namespace

	return &monitorV1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    m.Labels,
		},
		Spec: monitorV1.ServiceMonitorSpec{
			Endpoints: []monitorV1.Endpoint{
				{
					Interval: model.ServiceMonitorInterval,
					Port:     model.ServiceMonitorPort,
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					model.RoleName: name,
				},
			},
		},
	}, nil
}
