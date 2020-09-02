package svc

import (
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *SVC) GenerateService(name string, svcType string) *corev1.Service {
	namespace := s.Workload.GetNamespace()
	var clusterIP string

	if svcType == "Headless" {
		clusterIP = "None"
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    s.Labels,
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/port":   "http",
				"prometheus.io/path":   "/metrics",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:     model.ClientPort,
					Protocol: corev1.ProtocolTCP,
					Name:     "client",
				},
				{
					Name:     "server",
					Port:     model.ServerPort,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "leader-election",
					Port:     model.LeaderElectionPort,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Port:     model.ExporterPort,
					Protocol: corev1.ProtocolTCP,
					Name:     model.ExporterPortName,
				},
				{
					Port:     model.AgentPort,
					Protocol: corev1.ProtocolTCP,
					Name:     model.AgentPortName,
				},
			},
			ClusterIP: clusterIP,
			Selector:  s.Labels,
		},
	}
}
