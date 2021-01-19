package rollout

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"

	commonsts "github.com/ghostbaby/zookeeper-operator/controllers/workload/common/sts"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (r *Rollout) RollingUpgrades() error {

	var (
		zkUrl        string
		toUpgrade    []corev1.Pod
		actualPods   []string
		podsToDelete []corev1.Pod
		leaderNode   corev1.Pod
	)

	name := r.Workload.GetName()
	namespace := r.Workload.GetNamespace()
	deletingPods := make([]corev1.Pod, 0)
	healthyPods := make(map[string]corev1.Pod)
	deletedPods := []corev1.Pod{}

	//获取期望节点数
	expectReplica := *r.ExpectSts.Spec.Replicas
	//获取当前节点数
	actualReplica := *r.ActualSts.Spec.Replicas

	//获取当前节点信息
	sts, podList, err := utils.GetStatefulSetPods(r.Client, r.Workload, r.Labels, r.Log)
	if err != nil {
		return err
	}

	//如果当前节点数小于配置节点数，等待所有节点全部启动
	if len(podList.Items) < int(*r.ActualSts.Spec.Replicas) {
		r.Log.Info(
			"Some pods still need to be created/deleted.",
			"namespace", namespace, "statefulset_name", name,
			"expected_pods_num", r.Workload.Spec.Replicas, "actual_pods_num", len(podList.Items),
		)
		return nil
	}

	if commonsts.IsUpgradeStsResource(r.ExpectSts, r.ActualSts) {
		//更新sts配置
		if err := r.Client.Update(r.ExpectSts); err != nil {
			return err
		}
	}

	//如果当前节点数不等于期待节点数，重新执行主流程
	if actualReplica != expectReplica {
		return errors.New("rollingupdate need requeue.")
	}

	r.Log.Info(
		"Start to check ZooKeeper cluster resource.",
	)

	currentPods, err := utils.GetCurrentPods(r.Client, r.Workload, r.Labels, r.Log)
	if err != nil {
		r.Log.Info(
			"Unable to get current pods.",
			"error", err,
			"namespace", namespace,
			"zk_name", name,
		)
		return err
	}

	if len(currentPods) != int(actualReplica) {
		r.Log.Info(
			"Not all pod ready .",
			"error", err,
			"namespace", namespace,
			"zk_name", name,
		)
		return nil
	}

	r.Workload.Status.Phase = cachev1alpha1.ZooKeeperApplyingChangesPhase

	podArray := podList.Items

	zkUrl, _ = utils.GetServiceUrl(r.Workload, podArray)

	cli := &zk.BaseClient{
		HTTP:     &http.Client{},
		Endpoint: zkUrl,
	}

	for _, pod := range podArray {
		if pod.DeletionTimestamp != nil {
			deletingPods = append(deletingPods, pod)
			continue
		}
		if !pod.DeletionTimestamp.IsZero() || !utils.IsPodReady(pod) {
			continue
		}
		cli.Endpoint = fmt.Sprintf("%s://%s:%d", "http", pod.Status.PodIP, model.AgentHTTPPort)
		isLeader, inCluster, err := NodesInCluster(cli, []string{pod.Name})
		if err != nil {
			return err
		}
		if inCluster {
			healthyPods[pod.Name] = pod
		}
		if isLeader {
			leaderNode = pod
		}

		actualPods = append(actualPods, pod.Name)
		alreadyUpgraded := podUpgradeDone(pod, sts.Status.UpdateRevision)
		if !alreadyUpgraded && !isLeader {
			toUpgrade = append(toUpgrade, pod)
		}
	}

	//leader 节点放到最后升级
	leaderUpgraded := podUpgradeDone(leaderNode, sts.Status.UpdateRevision)
	if !leaderUpgraded {
		toUpgrade = append(toUpgrade, leaderNode)
	}

	if len(toUpgrade) == 0 {
		r.Log.Info(
			"No pod need to update.",
			"zk_name", name,
			"namespace", namespace,
		)
		return nil
	}

	//获取允许删除pod数量，以及maxUnavailable是否配置
	unhealthyPods := len(actualPods) - len(healthyPods)

	maxUnavailable := cachev1alpha1.DefaultChangeBudget.MaxUnavailable
	allowedDeletions := int(*maxUnavailable) - unhealthyPods

	maxUnavailableReached := allowedDeletions <= 0

	candidates := make([]corev1.Pod, len(toUpgrade))
	copy(candidates, toUpgrade)

	for _, candidate := range candidates {
		if len(deletingPods) > 0 {
			r.Log.Info("Allow delete quota is reached.")
			continue
		}
		//只允许一个健康节点进行重启
		_, healthy := healthyPods[candidate.Name]
		if maxUnavailableReached && healthy {
			r.Log.Info(
				"do_not_restart_healthy_node_if_MaxUnavailable_reached",
				"pod_name", candidate.Name,
				"zk_name", name,
				"namespace", namespace,
			)
			continue
		}

		//跳过状态为terminating的节点
		if candidate.DeletionTimestamp != nil {
			r.Log.Info(
				"skip_already_terminating_pods",
				"pod_name", candidate.Name,
				"zk_name", name,
				"namespace", namespace,
			)
			continue
		}

		delete(healthyPods, candidate.Name)
		podsToDelete = append(podsToDelete, candidate)
		allowedDeletions--
		if allowedDeletions <= 0 {
			r.Log.Info("Allow delete quota is reached.")
			break
		}

	}

	if len(podsToDelete) == 0 {
		r.Log.V(1).Info(
			"No pod deleted during rolling upgrade",
			"zk_name", name,
			"namespace", namespace,
		)
		return nil
	}

	for _, podToDelete := range podsToDelete {
		err := r.Client.Delete(&podToDelete)
		if err != nil {
			return err
		}
		deletedPods = append(deletedPods, podToDelete)
	}

	if len(deletedPods) > 0 {
		return errors.New("rollingupdate need requeue.")
	}

	if len(podsToDelete) > len(deletedPods) {
		return errors.New("rollingupdate need requeue.")
	}

	return nil
}

func NodesInCluster(cli *zk.BaseClient, nodeNames []string) (bool, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), zk.DefaultReqTimeout)
	defer cancel()
	nodes, err := cli.GetClusterStatus(ctx)
	if err != nil {
		return false, false, err
	}

	return nodes.Mode == 1, nodes.Mode != 0, nil
}

// podUpgradeDone inspects the given pod and returns true if it was successfully upgraded.
func podUpgradeDone(pod corev1.Pod, expectedRevision string) bool {
	if expectedRevision == "" {
		// no upgrade scheduled for the sset
		return false
	}
	if PodRevision(pod) != expectedRevision {
		// pod revision does not match the sset upgrade revision
		return false
	}
	return true
}

func PodRevision(pod corev1.Pod) string {
	return pod.Labels[appsv1.StatefulSetRevisionLabel]
}
