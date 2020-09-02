package cm

import (
	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
)

type CM struct {
	Workload  *cachev1alpha1.Workload
	Labels    map[string]string
	ExpectSts *appsv1.StatefulSet
	ActualSts *appsv1.StatefulSet
}

func NewCM(workload *cachev1alpha1.Workload, labels map[string]string, expectSTS, actualSTS *appsv1.StatefulSet) *CM {
	return &CM{
		Workload:  workload,
		Labels:    labels,
		ExpectSts: expectSTS,
		ActualSts: actualSTS,
	}
}
