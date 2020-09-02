package sts

import (
	"path"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/cm"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/volume"
	corev1 "k8s.io/api/core/v1"
)

var (
	// PodDownwardEnvVars inject the runtime Pod Name and IP as environment variables.
	PodDownwardEnvVars = []corev1.EnvVar{
		{Name: EnvPodIP, Value: "", ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "status.podIP"},
		}},
		{Name: EnvPodName, Value: "", ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"},
		}},
	}
)

func NewPrepareFSInitContainer(imageName string) (corev1.Container, error) {
	// we mount the certificates to a location outside of the default config directory because the prepare-fs script
	// will attempt to move all the files under the configuration directory to a different volume, and it should not
	// be attempting to move files from this secret volume mount (any attempt to do so will be logged as errors).

	scriptsVolume := corev1.VolumeMount{
		Name:      volume.ScriptsVolumeName,
		MountPath: volume.ScriptsVolumeMountPath,
		ReadOnly:  true,
	}

	privileged := false
	container := corev1.Container{
		Image:           imageName,
		ImagePullPolicy: corev1.PullAlways,
		Name:            PrepareFilesystemContainerName,
		SecurityContext: &corev1.SecurityContext{
			Privileged: &privileged,
		},
		Env:     PodDownwardEnvVars,
		Command: []string{"bash", "-c", path.Join(volume.ScriptsVolumeMountPath, cm.PrepareFsScriptConfigKey)},
		VolumeMounts: append(
			InitContainerVolumeMounts(cm.PluginVolumes),
			scriptsVolume,
			volume.DefaultDataVolumeMount,
			volume.DefaultLogsVolumeMount,
		),
	}

	return container, nil
}

func InitContainerVolumeMount(v cm.SharedVolume) corev1.VolumeMount {
	return corev1.VolumeMount{
		MountPath: v.InitContainerMountPath,
		Name:      v.Name,
	}
}

func InitContainerVolumeMounts(v cm.SharedVolumeArray) []corev1.VolumeMount {
	mounts := make([]corev1.VolumeMount, len(v.Array))
	for i, v := range v.Array {
		mounts[i] = InitContainerVolumeMount(v)
	}
	return mounts
}
