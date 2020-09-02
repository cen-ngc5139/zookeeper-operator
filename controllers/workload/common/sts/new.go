package sts

import (
	"github.com/ghostbaby/zookeeper-operator/controllers/common/volume"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/cm"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *STS) GenerateStatefulset() (*appsv1.StatefulSet, error) {
	name := s.Workload.GetName()
	namespace := s.Workload.GetNamespace()

	sts := NewStatefulSet(name, namespace)
	resource := getResources(s.Workload)

	var volumes []corev1.Volume
	var volumeMount []corev1.VolumeMount
	var image string
	var initContainerArray []corev1.Container

	if s.Workload.Spec.Image != "" {
		image = s.Workload.Spec.Image
	} else {
		image = utils.Joins(model.DefaultImageRepository, ":", s.Workload.Spec.Version)
	}

	//生成agent配置文件
	agentConfigFileVolume := volume.ConfigMapVolume{
		CmName:      genConfigMapName(name, cm.AgentVolumeName),
		Name:        cm.AgentVolumeName,
		MountPath:   cm.AgentVolumeMountPath,
		DefaultMode: 0755,
	}
	volumes = append(volumes, agentConfigFileVolume.Volume())
	volumeMount = append(volumeMount, agentConfigFileVolume.VolumeMount())

	//生成data存储卷
	dataVolume := GetDataVolume(s.Workload)
	if dataVolume != nil {
		volumes = append(volumes, *dataVolume)
	}
	containerDefaultVM := BuildDefaultVolumes()

	for _, vm := range containerDefaultVM {
		volumeMount = append(volumeMount, vm)
	}

	for _, volume := range cm.PluginVolumes.Volumes() {
		volumes = append(volumes, volume)
	}

	volumes = append(volumes, volume.DefaultLogsVolume)

	//生成初始化脚本存储卷
	scriptsVolume := volume.ConfigMapVolume{
		CmName:      genConfigMapName(name, cm.ScriptsVolumeName),
		Name:        cm.ScriptsVolumeName,
		MountPath:   cm.ScriptsVolumeMountPath,
		DefaultMode: 0755,
	}

	volumes = append(volumes, scriptsVolume.Volume())
	volumeMount = append(volumeMount, scriptsVolume.VolumeMount())

	//生成动态配置存储卷
	dynamicConfigFileVolume := volume.ConfigMapVolume{
		CmName:      genConfigMapName(name, cm.DynamicConfigFileVolumeName),
		Name:        cm.DynamicConfigFileVolumeName,
		MountPath:   cm.DynamicConfigFileVolumeMountPath,
		DefaultMode: 0755,
	}
	volumes = append(volumes, dynamicConfigFileVolume.Volume())
	volumeMount = append(volumeMount, dynamicConfigFileVolume.VolumeMount())

	//生成启动配置文件
	ConfigVolume := volume.ConfigMapVolume{
		CmName:      genConfigMapName(name, cm.ConfigVolumeName),
		Name:        cm.ConfigVolumeName,
		MountPath:   cm.ConfigVolumeMountPath,
		DefaultMode: 0755,
	}

	volumes = append(volumes, ConfigVolume.Volume())
	volumeMount = append(volumeMount, ConfigVolume.VolumeMount())

	//生成初始化容器
	prepareFsContainer, err := NewPrepareFSInitContainer(image)
	if err != nil {
		return nil, err
	}
	initContainerArray = append(initContainerArray, prepareFsContainer)

	sts.Spec = appsv1.StatefulSetSpec{
		ServiceName: name,
		Replicas:    &s.Workload.Spec.Cluster.NodeCount,
		Selector: &metav1.LabelSelector{
			MatchLabels: s.Labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      s.Labels,
				Annotations: s.Workload.Spec.Annotations,
			},
			Spec: corev1.PodSpec{
				Affinity:          PodAffinity(s.Workload.Spec.Affinity, s.Labels),
				NodeSelector:      s.Workload.Spec.NodeSelector,
				Tolerations:       s.Workload.Spec.Tolerations,
				PriorityClassName: s.Workload.Spec.PriorityClassName,
				RestartPolicy:     corev1.RestartPolicyAlways,
				InitContainers:    initContainerArray,
				Containers: []corev1.Container{
					container(resource, image, volumeMount),
				},
				Volumes: volumes,
			},
		},
		UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.OnDeleteStatefulSetStrategyType,
		},
		PodManagementPolicy: appsv1.ParallelPodManagement,
	}

	if s.Workload.Spec.Cluster.Exporter.Exporter {
		exporter := s.CreateExporterContainer(s.Workload)
		sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, exporter)
	}

	//生成agent容器
	agentContainer := s.CreateAgentContainer(agentConfigFileVolume.VolumeMount())
	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, agentContainer)

	if s.Workload.Spec.Cluster.Storage.PersistentVolumeClaim != nil {
		sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
			*s.Workload.Spec.Cluster.Storage.PersistentVolumeClaim,
		}
	}

	return sts, nil
}
