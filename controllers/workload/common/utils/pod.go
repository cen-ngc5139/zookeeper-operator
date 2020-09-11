package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/labels"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodName returns the name of the pod with the given ordinal for this StatefulSet.
func PodName(ssetName string, ordinal int32) string {
	return fmt.Sprintf("%s-%d", ssetName, ordinal)
}

func genPodHost(rf *cachev1alpha1.Workload, podName string) string {
	return fmt.Sprintf("%s.%s", podName, getDomain(rf))
}

func getDomain(rf *cachev1alpha1.Workload) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", getSvcName(rf), rf.GetNamespace())
}

func getSvcName(rf *cachev1alpha1.Workload) string {
	return rf.GetName()
}

func PodNames(sset appsv1.StatefulSet) []string {
	names := make([]string, 0, GetReplicas(sset))
	for i := int32(0); i < GetReplicas(sset); i++ {
		names = append(names, PodName(sset.Name, i))
	}
	return names
}

func GetReplicas(sts appsv1.StatefulSet) int32 {
	if sts.Spec.Replicas != nil {
		return *sts.Spec.Replicas
	}
	return 0
}

func GetByName(ssetName string, expectName string) bool {
	return expectName == ssetName
}

// ExtractNamespacedName returns an NamespacedName based on the given Object.
func ExtractNamespacedName(object metav1.Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: object.GetNamespace(),
		Name:      object.GetName(),
	}
}

// IsAvailable checks if both conditions ContainersReady and PodReady of a Pod are true.
func IsPodReady(pod corev1.Pod) bool {
	conditionsTrue := 0
	for _, cond := range pod.Status.Conditions {
		if cond.Status == corev1.ConditionTrue && (cond.Type == corev1.ContainersReady || cond.Type == corev1.PodReady) {
			conditionsTrue++
		}
	}
	return conditionsTrue == 2
}

func GetStatefulSetPods(c k8s.Client, w *cachev1alpha1.Workload, label map[string]string, log logr.Logger) (*appsv1.StatefulSet, *corev1.PodList, error) {
	sts := &appsv1.StatefulSet{}
	name := w.GetName()
	namespace := w.GetNamespace()
	err := c.Get(types.NamespacedName{Name: name, Namespace: namespace}, sts)
	if err != nil {
		return nil, nil, err
	}

	opts := &client.ListOptions{}
	set := labels.SelectorFromSet(label)
	opts.LabelSelector = set

	pod := &corev1.PodList{}

	if err := c.List(opts, pod); err != nil {
		log.Error(err, "fail to get pod.", "namespace", namespace, "name", name)
		return nil, nil, err
	}
	return sts, pod, nil
}

func GetCurrentPods(c k8s.Client, w *cachev1alpha1.Workload, label map[string]string, log logr.Logger) ([]corev1.Pod, error) {
	_, podList, err := GetStatefulSetPods(c, w, label, log)
	if err != nil {
		return nil, err
	}
	currentPods := make([]corev1.Pod, 0, len(podList.Items))

	for _, p := range podList.Items {
		var containerStatusReday bool
		var isNewPod bool
		if p.DeletionTimestamp != nil {
			continue
		}
		containerStatusReday, err = PodRunningAndReady(p)
		if err != nil {
			return nil, err
		}

		if !containerStatusReday && !isNewPod {
			continue
		}
		currentPods = append(currentPods, p)
	}
	return currentPods, nil
}

// PodRunningAndReady returns whether a pod is running and each container has
// passed it's ready state.
func PodRunningAndReady(pod corev1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case corev1.PodFailed, corev1.PodSucceeded:
		return false, fmt.Errorf("pod completed")
	case corev1.PodRunning:
		for _, cond := range pod.Status.Conditions {
			if cond.Type != corev1.PodReady {
				continue
			}
			return cond.Status == corev1.ConditionTrue, nil
		}
		return false, fmt.Errorf("pod ready condition not found")
	}
	return false, nil
}

// getReplsetAddrs returns a slice of replset host:port addresses
func GetPodIp(w *cachev1alpha1.Workload, podNames []string) []string {
	addrs := make([]string, 0)
	for _, podName := range podNames {

		podIndexArray := strings.Split(podName, "-")
		podIndex, _ := strconv.Atoi(podIndexArray[len(podIndexArray)-1])
		addrs = append(
			addrs,
			fmt.Sprintf("server.%d=%s:%d:%d:participant;0.0.0.0:%d",
				podIndex+1, genPodHost(w, podName), model.ServerPort, model.LeaderElectionPort, model.ClientPort),
		)
	}
	return addrs
}

func GetServiceUrl(w *cachev1alpha1.Workload, pods []corev1.Pod) (string, []string) {
	var urls []string
	var url string
	_, err := rest.InClusterConfig()
	if err != nil {
		randomPod := pods[rand.Intn(len(pods))]
		for _, pod := range pods {
			url := fmt.Sprintf("%s://%s:%d", "http", pod.Status.PodIP, model.AgentHTTPPort)
			urls = append(urls, url)
		}
		url = fmt.Sprintf("%s://%s:%d", "http", randomPod.Status.PodIP, model.AgentHTTPPort)
		return url, urls
	}

	for _, pod := range pods {
		url := fmt.Sprintf("%s://%s:%d", "http", genPodHost(w, pod.Name), model.AgentHTTPPort)
		urls = append(urls, url)
	}

	randomPod := pods[rand.Intn(len(pods))]
	url = fmt.Sprintf("%s://%s:%d", "http", genPodHost(w, randomPod.Name), model.AgentHTTPPort)
	return url, urls
}
