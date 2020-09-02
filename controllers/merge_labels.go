package controllers

import (
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
)

func GenerateLabels(labels map[string]string, name string) map[string]string {
	dynLabels := map[string]string{
		model.AppLabel:                 name,
		"app.kubernetes.io/name":       "zookeeper",
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/managed-by": "zookeeper-operator",
		"app.kubernetes.io/part-of":    "zookeeper",
	}
	return MergeLabels(dynLabels, labels)
}

func MergeLabels(allLabels ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, labels := range allLabels {
		if labels != nil {
			for k, v := range labels {
				res[k] = v
			}
		}
	}
	return res
}
