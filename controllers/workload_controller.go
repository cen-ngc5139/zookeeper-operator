/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/prometheus"

	cachev1alpha1 "github.com/ghostbaby/zookeeper-operator/api/v1alpha1"
	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/finalizer"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/observer"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/model"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/provision"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// WorkloadReconciler reconciles a Workload object
type WorkloadReconciler struct {
	client.Client
	ServiceGetter
	Log           logr.Logger
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	Observers     *observer.Manager
	Monitor       *prometheus.GenericClientset
	ZKClient      *zk.BaseClient
	ObservedState *observer.State
	Finalizers    finalizer.Handler
}

var ReconcileWaitResult = reconcile.Result{RequeueAfter: 30 * time.Second}

// +kubebuilder:rbac:groups=cache.ghostbaby.io,resources=workloads,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.ghostbaby.io,resources=workloads/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods;configmaps;services;events;secret,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors;prometheusrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods;configmaps;services;events;secret,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/exec,verbs=create

func (r *WorkloadReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("workload", req.NamespacedName)

	log.Info("start to reconcile.")

	var workload cachev1alpha1.Workload
	if err := r.Get(ctx, req.NamespacedName, &workload); err != nil {
		log.Error(err, "unable to fetch Zookeeper Workload")
		return ReconcileWaitResult, client.IgnoreNotFound(err)
	}

	//// workload will be gracefully deleted by server when DeletionTimestamp is non-null
	//if workload.DeletionTimestamp != nil {
	//	return ReconcileWaitResult, nil
	//}

	dClient, err := k8s.NewDynamicClient()
	if err != nil {
		log.Error(err, "unable to create dynamic client")
		return ReconcileWaitResult, err
	}

	option := &GetOptions{
		Client:        k8s.WrapClient(ctx, r.Client),
		Recorder:      r.Recorder,
		Log:           r.Log,
		DClient:       k8s.WrapDClient(dClient),
		Scheme:        r.Scheme,
		Labels:        GenerateLabels(workload.Labels, workload.Name),
		Observers:     r.Observers,
		ZKClient:      r.ZKClient,
		ObservedState: r.ObservedState,
		Finalizers:    r.Finalizers,
		Monitor:       r.Monitor,
	}

	if err := r.Workload(ctx, &workload, option).Reconcile(); err != nil {
		log.Error(err, "error when reconcile workload.")
		return ReconcileWaitResult, err
	}

	return ctrl.Result{}, nil
}

func (r *WorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mapFn := handler.ToRequestsFunc(
		func(a handler.MapObject) []reconcile.Request {

			labels := a.Meta.GetLabels()
			clusterName, isSet := labels[model.RoleName]
			if !isSet {
				return nil
			}
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      clusterName,
					Namespace: a.Meta.GetNamespace(),
				}},
			}
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Workload{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Pod{}).
		Watches(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: mapFn}).
		Watches(observer.WatchClusterHealthChange(r.Observers), provision.GenericEventHandler()).
		Complete(r)
}
