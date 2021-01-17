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

package v1alpha1

import (
	"strconv"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var workloadlog = logf.Log.WithName("workload-resource")

func (r *Workload) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-zk-cache-ghostbaby-io-v1alpha1-workload,mutating=true,failurePolicy=fail,groups=cache.ghostbaby.io,resources=workloads,verbs=create;update,versions=v1alpha1,name=mworkload.kb.io

var _ webhook.Defaulter = &Workload{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Workload) Default() {
	workloadlog.Info("default", "name", r.Name)

	if r.Spec.Cluster.Resources == nil {
		r.Spec.Cluster.Resources = &Resources{
			Requests: CPUAndMem{
				CPU:    "1000m",
				Memory: "3Gi",
			},
			Limits: CPUAndMem{
				CPU:    "2000m",
				Memory: "4Gi",
			},
		}
	}

	if r.Spec.Cluster.Exporter == nil {
		r.Spec.Cluster.Exporter = &ExporterSpec{
			Exporter:              true,
			ExporterImage:         "ghostbaby/zookeeper_exporter",
			ExporterVersion:       "v3.5.6",
			DisableExporterProbes: false,
		}
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-zk-cache-ghostbaby-io-v1alpha1-workload,mutating=false,failurePolicy=fail,groups=cache.ghostbaby.io,resources=workloads,versions=v1alpha1,name=vworkload.kb.io

var _ webhook.Validator = &Workload{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Workload) ValidateCreate() error {
	workloadlog.Info("validate create", "name", r.Name)

	return r.validateZooKeeper()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Workload) ValidateUpdate(old runtime.Object) error {
	workloadlog.Info("validate update", "name", r.Name)

	return r.validateZooKeeper()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Workload) ValidateDelete() error {
	workloadlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Workload) validateZooKeeper() error {
	var allErrs field.ErrorList

	if err := r.CheckResourceCpu(
		r.Spec.Cluster.Resources.Limits.CPU,
		r.Spec.Cluster.Resources.Requests.CPU,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.CheckResourceMem(
		r.Spec.Cluster.Resources.Limits.Memory,
		r.Spec.Cluster.Resources.Requests.Memory,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "zk.cache.ghostbaby.io", Kind: "Workload"},
		r.Name, allErrs)
}

func (r *Workload) CheckResourceMem(limit string, request string) *field.Error {

	limitMem := ResourceMemString2Int(limit)
	requestMem := ResourceMemString2Int(request)
	if limitMem < requestMem {
		return field.Invalid(
			field.NewPath("resources").Child("limit"),
			r.Name, "limit mem must greater than request mem.")
	}

	return nil
}

func (r *Workload) CheckResourceCpu(limit string, request string) *field.Error {

	limitCPU := ResourceCpuString2Int(limit)
	requestCPU := ResourceCpuString2Int(request)
	if limitCPU < requestCPU {
		return field.Invalid(
			field.NewPath("resources").Child("limit"),
			r.Name, "limit cpu must greater than request cpu.")
	}

	return nil

}

func ResourceMemString2Int(resource string) int64 {
	var (
		resourceStr string
		resourceInt int64
	)
	if strings.Contains(resource, "Gi") {
		resourceStr = strings.Split(resource, "Gi")[0]
		resourceInt, _ = strconv.ParseInt(resourceStr, 10, 32)
		resourceInt = resourceInt * 1024
	} else if strings.Contains(resource, "Mi") {
		resourceStr = strings.Split(resource, "Mi")[0]
		resourceInt, _ = strconv.ParseInt(resourceStr, 10, 32)
	}
	return resourceInt
}

func ResourceCpuString2Int(resource string) int64 {
	var (
		resourceStr string
		resourceInt int64
	)
	if strings.Contains(resource, "m") {
		resourceStr = strings.Split(resource, "m")[0]
		resourceInt, _ = strconv.ParseInt(resourceStr, 10, 32)
		resourceInt = resourceInt / 1000
	} else {

		resourceInt, _ = strconv.ParseInt(resource, 10, 32)
	}
	return resourceInt
}
