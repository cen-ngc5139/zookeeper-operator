package provision

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"

	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	DefaultObservationInterval = 10 * time.Second
	DefaultRequestTimeout      = 1 * time.Minute
)

// ResourcesState contains information about a deployments resources.
type ResourcesState struct {
	// AllPods are all the pods related to the cluster, including ones with a
	// DeletionTimestamp tombstone set.
	AllPods []corev1.Pod
	// CurrentPods are all non-deleted pods.
	CurrentPods []corev1.Pod
	// CurrentPodsByPhase are all non-deleted indexed by their PodPhase
	CurrentPodsByPhase map[corev1.PodPhase][]corev1.Pod
	// DeletingPods are all deleted pods.
	DeletingPods []corev1.Pod
}

func (p *Provision) Observer() error {
	var zkUrls []string
	var zkUrl string

	_, podList, err := utils.GetStatefulSetPods(p.Client, p.Workload, p.Labels, p.Log)
	if err != nil {
		return err
	}

	if len(podList.Items) == 0 {
		p.Log.Info("pod list is empty，pls wait.")
		return nil
	}

	podArray := podList.Items

	deletingPods := make([]corev1.Pod, 0)
	currentPods := make([]corev1.Pod, 0, len(podArray))
	currentPodsByPhase := make(map[corev1.PodPhase][]corev1.Pod)

	for _, p := range podArray {
		if p.DeletionTimestamp != nil {
			deletingPods = append(deletingPods, p)
			continue
		}
		currentPods = append(currentPods, p)
		podsInPhase, ok := currentPodsByPhase[p.Status.Phase]
		if !ok {
			podsInPhase = []corev1.Pod{p}
		} else {
			podsInPhase = append(podsInPhase, p)
		}
		currentPodsByPhase[p.Status.Phase] = podsInPhase
	}

	//podState := ResourcesState{
	//	AllPods:            podArray,
	//	CurrentPods:        currentPods,
	//	CurrentPodsByPhase: currentPodsByPhase,
	//	DeletingPods:       deletingPods,
	//}

	zkUrl, zkUrls = utils.GetServiceUrl(p.Workload, currentPods[0:2])

	cli := &zk.BaseClient{
		Endpoints: zkUrls,
		HTTP:      &http.Client{},
		Endpoint:  zkUrl,
		Transport: &http.Transport{},
	}

	state := p.Observers.ObservedStateResolver(
		utils.ExtractNamespacedName(p.Workload),
		*cli,
	)

	if state.ClusterStats != nil {
		p.Workload.Status.LeaderNode = state.ClusterStats.LeaderNode
		p.Workload.Status.AvailableNodes = state.ClusterStats.AvailableNodes
		p.Workload.Status.LastTransitionTime = metav1.Now()
	}

	p.ObservedState = &state
	p.ZKClient = cli

	p.writeStatus()

	return nil
}

func (p *Provision) writeStatus() error {
	err := p.Client.WriteStatus(p.Workload)
	if err != nil {
		// may be it's k8s v1.10 and erlier (e.g. oc3.9) that doesn't support status updates
		// so try to update whole CR
		err := p.Client.Update(p.Workload)
		if err != nil {
			return errors.Wrap(err, "send update")
		}
	}

	return nil
}

// GenericEventHandler returns an EventHandler that enqueues a reconciliation request
// from the generic event NamespacedName.
func GenericEventHandler() handler.EventHandler {
	return handler.Funcs{
		GenericFunc: func(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
			q.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: evt.Meta.GetNamespace(),
					Name:      evt.Meta.GetName(),
				},
			})
		},
	}
}
