package svc

import cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"

type SVC struct {
	Workload *cachev1alpha1.Workload
	Labels   map[string]string
}

func NewSVC(workload *cachev1alpha1.Workload, labels map[string]string) *SVC {
	return &SVC{
		Workload: workload,
		Labels:   labels,
	}
}
