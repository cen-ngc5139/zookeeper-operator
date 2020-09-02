package sts

import cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"

type STS struct {
	Workload *cachev1alpha1.Workload
	Labels   map[string]string
}

func NewSTS(workload *cachev1alpha1.Workload, labels map[string]string) *STS {
	return &STS{
		Workload: workload,
		Labels:   labels,
	}
}
