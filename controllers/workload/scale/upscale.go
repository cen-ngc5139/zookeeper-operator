package scale

import (
	"fmt"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (s *Scale) UpScale() error {
	//[BUG] upscale.go:13 +0x51 集群刚刚启动，无法获取到期待sts，导致operator crush

	if s.ExpectSts == nil || s.ActualSts == nil {
		return nil
	}
	name := s.Workload.GetName()
	expectReplica := s.ExpectSts.Spec.Replicas
	actualReplica := s.ActualSts.Spec.Replicas

	if expectReplica != nil && actualReplica != nil && *expectReplica > *actualReplica {
		msg := fmt.Sprintf(model.UpdateMessageZooKeeperStatefulset, name)
		s.Recorder.Event(s.Workload, corev1.EventTypeNormal, model.ZooKeeperStatefulset, msg)

		s.Log.Info(
			"Scaling replicas up",
			"from", actualReplica,
			"to", expectReplica,
		)

		s.Workload.Status.Phase = cachev1alpha1.ZooKeeperUpScaling
		msg = fmt.Sprintf(model.MessageZooKeeperUpScaling, actualReplica, expectReplica)
		s.Recorder.Event(s.Workload, corev1.EventTypeNormal, model.ZooKeeperUpScaling, msg)

		err := s.Client.Update(s.ExpectSts)
		if err != nil {
			return err
		}
	}
	return nil
}
