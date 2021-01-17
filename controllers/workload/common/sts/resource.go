package sts

import (
	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func getResources(w *cachev1alpha1.Workload) corev1.ResourceRequirements {

	reqCpu := w.Spec.Cluster.Resources.Requests.CPU
	reqMem := w.Spec.Cluster.Resources.Requests.Memory
	limCpu := w.Spec.Cluster.Resources.Limits.CPU
	limMem := w.Spec.Cluster.Resources.Limits.Memory

	return corev1.ResourceRequirements{
		Requests: getRequests(reqCpu, reqMem),
		Limits:   getLimits(limCpu, limMem),
	}
}

func getLimits(limCpu string, limitMem string) corev1.ResourceList {
	return generateResourceList(limCpu, limitMem)
}

func getRequests(reqCpu string, reqMem string) corev1.ResourceList {
	return generateResourceList(reqCpu, reqMem)
}

func generateResourceList(cpu string, memory string) corev1.ResourceList {
	resources := corev1.ResourceList{}
	if cpu != "" {
		resources[corev1.ResourceCPU], _ = resource.ParseQuantity(cpu)
	}
	if memory != "" {
		resources[corev1.ResourceMemory], _ = resource.ParseQuantity(memory)
	}
	return resources
}

func getStsResource(sts *appsv1.StatefulSet) *cachev1alpha1.ZkResource {
	var requestCpu, requestMem *resource.Quantity
	var limitCpu, limitMem *resource.Quantity
	for _, container := range sts.Spec.Template.Spec.Containers {
		if container.Name == model.RoleName {
			requestCpu = container.Resources.Requests.Cpu()
			requestMem = container.Resources.Requests.Memory()

			limitCpu = container.Resources.Limits.Cpu()
			limitMem = container.Resources.Limits.Memory()
		} else {
			continue
		}
	}
	return &cachev1alpha1.ZkResource{
		RequestCpu: requestCpu,
		RequestMem: requestMem,
		LimitCpu:   limitCpu,
		LimitMem:   limitMem,
	}
}

func IsUpgradeStsResource(expectSts *appsv1.StatefulSet, actualSts *appsv1.StatefulSet) bool {
	expectEsResource := getStsResource(expectSts)
	actualEsResource := getStsResource(actualSts)
	//fmt.Println(expectEsResource, actualEsResource)
	if expectEsResource.RequestCpu.String() != actualEsResource.RequestCpu.String() ||
		expectEsResource.RequestMem.String() != actualEsResource.RequestMem.String() ||
		expectEsResource.LimitCpu.String() != actualEsResource.LimitCpu.String() ||
		expectEsResource.LimitMem.String() != actualEsResource.LimitMem.String() {
		return true
	} else {
		return false
	}
}
