package volume

import (
	corev1 "k8s.io/api/core/v1"
)

type ConfigMapVolume struct {
	CmName      string
	Name        string
	MountPath   string
	DefaultMode int32
}

var (
	defaultOptional = false
)

func (cm ConfigMapVolume) Volume() corev1.Volume {
	return corev1.Volume{
		Name: cm.Name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cm.CmName,
				},
				Optional:    &defaultOptional,
				DefaultMode: &cm.DefaultMode,
			},
		},
	}
}

// VolumeMount returns the k8s volume mount.
func (cm ConfigMapVolume) VolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      cm.Name,
		MountPath: cm.MountPath,
		ReadOnly:  true,
	}
}
