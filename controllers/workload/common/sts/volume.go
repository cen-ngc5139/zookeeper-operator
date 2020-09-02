package sts

import (
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/volume"
	corev1 "k8s.io/api/core/v1"
)

func BuildDefaultVolumes() []corev1.VolumeMount {
	var volumeMounts []corev1.VolumeMount
	volumeMounts = append(volumeMounts,
		volume.DefaultDataVolumeMount,
		volume.DefaultLogsVolumeMount,
	)

	return volumeMounts
}
