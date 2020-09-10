package sts

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"

	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/volume"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewStatefulSet returns a StatefulSet object configured for a name
func NewStatefulSet(name, namespace string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}

// PodAffinity returns podAffinity options for the pod
func PodAffinity(af *cachev1alpha1.PodAffinity, labels map[string]string) *corev1.Affinity {
	if af == nil {
		return nil
	}

	switch {
	case af.Advanced != nil:
		return af.Advanced
	case af.TopologyKey != nil:
		if *af.TopologyKey == volume.AffinityOff {
			return nil
		}

		l := make(map[string]string)
		for k, v := range labels {
			if k != "app.kubernetes.io/component" {
				l[k] = v
			}
		}
		return &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: l,
						},
						TopologyKey: *af.TopologyKey,
					},
				},
			},
		}
	}

	return nil
}

func GetExporterImage(rf *cachev1alpha1.Workload) string {
	return fmt.Sprintf("%s:%s", rf.Spec.Cluster.Exporter.ExporterImage, rf.Spec.Cluster.Exporter.ExporterVersion)
}

func (r *STS) CreateExporterContainer(rf *cachev1alpha1.Workload) corev1.Container {

	exporterImage := GetExporterImage(rf)

	// Define readiness and liveness probes only if config option to disable isn't set
	var readinessProbe, livenessProbe *corev1.Probe
	if !rf.Spec.Cluster.Exporter.DisableExporterProbes {
		readinessProbe = &corev1.Probe{
			InitialDelaySeconds: 10,
			TimeoutSeconds:      3,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/metrics",
					Port: intstr.FromString("metrics"),
				},
			},
		}

		livenessProbe = &corev1.Probe{
			TimeoutSeconds: 3,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/metrics",
					Port: intstr.FromString("metrics"),
				},
			},
		}
	}

	return corev1.Container{
		Name:            model.ExporterContainerName,
		Image:           exporterImage,
		ImagePullPolicy: "Always",
		Ports: []corev1.ContainerPort{
			{
				Name:          "metrics",
				ContainerPort: model.ExporterPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Command: []string{
			"sh",
			"-c",
			"/usr/local/bin/zookeeper-exporter  -listen 0.0.0.0:9114",
		},
		ReadinessProbe: readinessProbe,
		LivenessProbe:  livenessProbe,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(model.ExporterDefaultLimitCPU),
				corev1.ResourceMemory: resource.MustParse(model.ExporterDefaultLimitMemory),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(model.ExporterDefaultRequestCPU),
				corev1.ResourceMemory: resource.MustParse(model.ExporterDefaultRequestMemory),
			},
		},
	}
}

//生成jvm大小，如果小于1G，默认1G
func GenJvmHeapSize(mem string) (string, string) {
	var serverJvmHeapSize string
	var clientJvmHeapSize string
	if strings.Contains(mem, "Gi") {
		mem = strings.Split(mem, "Gi")[0]
		memInt32, _ := strconv.ParseInt(mem, 10, 32)

		if memInt32 > 0 {
			serverJvmHeapSize = fmt.Sprintf("%d", memInt32*1024)
			clientJvmHeapSize = fmt.Sprintf("%d", (memInt32*1024)/4)
		} else {
			serverJvmHeapSize = "1024"
			clientJvmHeapSize = "256"
		}
	}
	return serverJvmHeapSize, clientJvmHeapSize
}

func container(resources corev1.ResourceRequirements, image string, vm []corev1.VolumeMount) corev1.Container {

	runAsUser := int64(1000)

	serverJvmHeapSize, clientJvmHeapSize := GenJvmHeapSize(resources.Requests.Memory().String())

	return corev1.Container{
		Name:            model.ContainerName,
		Image:           image,
		ImagePullPolicy: corev1.PullAlways,
		Ports: []corev1.ContainerPort{
			{
				Name:          "client",
				ContainerPort: model.ClientPort,
			},
			{
				Name:          "server",
				ContainerPort: model.ServerPort,
			},
			{
				Name:          "leader-election",
				ContainerPort: model.LeaderElectionPort,
			},
		},
		Env: []corev1.EnvVar{
			{
				Name: "POD_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
			{
				Name:  "ZK_SERVER_HEAP",
				Value: serverJvmHeapSize,
			},
			{
				Name:  "ZK_CLIENT_HEAP",
				Value: clientJvmHeapSize,
			},
		},
		Command: []string{
			"sh",
			"-c",
			"zkStart.sh",
		},

		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  &runAsUser,
			RunAsGroup: &runAsUser,
		},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: 10,
			TimeoutSeconds:      10,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{Command: []string{"zkOK.sh"}},
			},
		},
		LivenessProbe: &corev1.Probe{
			InitialDelaySeconds: 10,
			TimeoutSeconds:      10,
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{Command: []string{"zkOK.sh"}},
			},
		},
		Resources:    resources,
		VolumeMounts: vm,
	}
}

func GetDataVolume(rf *cachev1alpha1.Workload) *corev1.Volume {
	// This will find the volumed desired by the user. If no volume defined
	// an EmptyDir will be used by default
	switch {
	case rf.Spec.Cluster.Storage.PersistentVolumeClaim != nil:
		return &corev1.Volume{
			Name: volume.DataVolClaimName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					// actual claim name will be resolved and fixed right before pod creation
					ClaimName: "claim-name-placeholder",
				},
			},
		}
	case rf.Spec.Cluster.Storage.EmptyDir != nil:
		return &corev1.Volume{
			Name: volume.DataVolClaimName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: rf.Spec.Cluster.Storage.EmptyDir,
			},
		}
	default:
		return &corev1.Volume{
			Name: volume.DataVolClaimName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
	}
}

func (r *STS) CreateAgentContainer(vm corev1.VolumeMount) corev1.Container {

	Image := "ghostbaby/zk-agent:v0.0.1"

	// Define readiness and liveness probes only if config option to disable isn't set
	var readinessProbe, livenessProbe *corev1.Probe

	readinessProbe = &corev1.Probe{
		InitialDelaySeconds: 10,
		TimeoutSeconds:      3,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.FromString("agent"),
			},
		},
	}

	livenessProbe = &corev1.Probe{
		TimeoutSeconds: 3,
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.FromString("agent"),
			},
		},
	}

	return corev1.Container{
		Name:            model.AgentContainerName,
		Image:           Image,
		ImagePullPolicy: "Always",
		Ports: []corev1.ContainerPort{
			{
				Name:          "agent",
				ContainerPort: model.AgentPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		Command: []string{
			"sh",
			"-c",
			"/usr/local/bin/zk-agent  -c /mnt/zookeeper/agent/config.json",
		},
		ReadinessProbe: readinessProbe,
		LivenessProbe:  livenessProbe,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(model.AgentDefaultLimitCPU),
				corev1.ResourceMemory: resource.MustParse(model.AgentDefaultLimitMemory),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(model.AgentDefaultRequestCPU),
				corev1.ResourceMemory: resource.MustParse(model.AgentDefaultRequestMemory),
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			vm,
		},
	}
}

func genConfigMapName(name, cmType string) string {
	return fmt.Sprintf("%s-%s", name, cmType)
}

func GetStatefulSet(c k8s.Client, w *cachev1alpha1.Workload, label map[string]string, scheme *runtime.Scheme,
) (*appsv1.StatefulSet, *appsv1.StatefulSet, error) {

	var (
		expectSts *appsv1.StatefulSet
		actualSts *appsv1.StatefulSet
	)

	actual := &appsv1.StatefulSet{}
	name := w.Name
	namespace := w.Namespace
	newSts := NewSTS(w, label)
	expect, err := newSts.GenerateStatefulset()
	if err != nil {
		return nil, nil, err
	}
	if err := controllerutil.SetControllerReference(w, expect, scheme); err != nil {
		return nil, nil, err
	}

	if err := c.Get(types.NamespacedName{Name: name, Namespace: namespace}, actual); err != nil && errors.IsNotFound(err) {
		expectSts = expect
	} else if err != nil {
		return nil, nil, err
	} else {
		expectSts = expect
		actualSts = actual
	}

	return expectSts, actualSts, nil
}
