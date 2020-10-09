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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WorkloadSpec defines the desired state of Workload
type WorkloadSpec struct {
	// Version represents the version of the stack
	Version string `json:"version,omitempty"`

	// Image represents the docker image that will be used.
	Image    string      `json:"image,omitempty"`
	Cluster  ClusterSpec `json:"cluster,omitempty"`
	Replicas *int32      `json:"replicas,omitempty"`

	//PodDisruptionBudget *PodDisruptionBudgetTemplate `json:"podDisruptionBudget,omitempty"`
	UpdateStrategy    UpdateStrategy      `json:"updateStrategy,omitempty"`
	Affinity          *PodAffinity        `json:"affinity,omitempty"`
	NodeSelector      map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations       []corev1.Toleration `json:"tolerations,omitempty"`
	PriorityClassName string              `json:"priorityClassName,omitempty"`
	Annotations       map[string]string   `json:"annotations,omitempty"`
	Labels            map[string]string   `json:"labels,omitempty"`
}

// GroupingDefinition is used to select a group of pods.
type GroupingDefinition struct {
	// Selector is the selector used to match pods.
	Selector metav1.LabelSelector `json:"selector,omitempty"`
}

// ChangeBudget defines how Pods in a single group should be updated.
type ChangeBudget struct {
	MaxUnavailable *int32 `json:"maxUnavailable"`

	MaxSurge *int32 `json:"maxSurge"`
}

type UpdateStrategy struct {
	// Groups is a list of groups that should have their cluster mutations considered in a fair manner with a strict
	// change budget (not allowing any surge or unavailability) before the entire cluster is reconciled with the
	// full change budget.
	Groups []GroupingDefinition `json:"groups,omitempty"`

	// ChangeBudget is the change budget that should be used when performing mutations to the cluster.
	ChangeBudget *ChangeBudget `json:"changeBudget,omitempty"`
}

// DefaultChangeBudget is used when no change budget is provided. It might not be the most effective, but should work in
// most cases.
var DefaultChangeBudget = ChangeBudget{
	MaxSurge:       nil,
	MaxUnavailable: Int32(1),
}

func Int32(v int32) *int32 { return &v }

type ClusterSpec struct {
	// Name is a logical name for this set of nodes. Used as a part of the managed Elasticsearch node.name setting.
	// +kubebuilder:validation:Pattern=[a-zA-Z0-9-]+
	// +kubebuilder:validation:MaxLength=23
	Name string `json:"name"`

	Storage Storage `json:"storage,omitempty"`

	Exporter  *ExporterSpec `json:"exporter,omitempty"`
	Resources *Resources    `json:"resources,omitempty"`
}

type ExporterSpec struct {
	Exporter              bool   `json:"exporter,omitempty"`
	ExporterImage         string `json:"exporterImage,omitempty"`
	ExporterVersion       string `json:"exporterVersion,omitempty"`
	DisableExporterProbes bool   `json:"disableExporterProbes,omitempty"`
}

// RedisResources sets the limits and requests for a container
type Resources struct {
	Requests CPUAndMem `json:"requests,omitempty"`
	Limits   CPUAndMem `json:"limits,omitempty"`
}

// CPUAndMem defines how many cpu and ram the container will request/limit
type CPUAndMem struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// RedisStorage defines the structure used to store the Redis Data
type Storage struct {
	KeepAfterDeletion     bool                          `json:"keepAfterDeletion,omitempty"`
	EmptyDir              *corev1.EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *corev1.PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`
	PersistentVolumeSize  string                        `json:"persistentVolumeSize,omitempty"`
}

type PodAffinity struct {
	TopologyKey *string          `json:"antiAffinityTopologyKey,omitempty"`
	Advanced    *corev1.Affinity `json:"advanced,omitempty"`
}

// PodDisruptionBudgetTemplate contains a template for creating a PodDisruptionBudget.
//type PodDisruptionBudgetTemplate struct {
//	// ObjectMeta is metadata for the service.
//	// The name and namespace provided here is managed by ECK and will be ignored.
//	// +optional
//	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`
//
//	// Spec of the desired behavior of the PodDisruptionBudget
//	// +optional
//	Spec policyv1beta1.PodDisruptionBudgetSpec `json:"spec,omitempty"`
//}

// WorkloadStatus defines the observed state of Workload
type WorkloadStatus struct {
	LeaderNode string `json:"leaderNode,omitempty"`

	// ReadyReplicas is the number of number of ready replicas in the cluster
	AvailableNodes int `json:"availableNodes,omitempty"`

	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	Phase ZooKeeperOrchestrationPhase `json:"phase,omitempty"`

	// The generation observed by the appConfig controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`
}

// MembersStatus is the status of the members of the cluster with both
// ready and unready node membership lists
type MembersStatus struct {
	Ready   []string `json:"ready"`
	Unready []string `json:"unready"`
}

type ZkResource struct {
	RequestCpu *resource.Quantity
	RequestMem *resource.Quantity
	LimitCpu   *resource.Quantity
	LimitMem   *resource.Quantity
}

type ZooKeeperOrchestrationPhase string

const (
	// ElasticsearchReadyPhase is operating at the desired spec.
	ZooKeeperReadyPhase ZooKeeperOrchestrationPhase = "Ready"
	// ZooKeeperApplyingChangesPhase controller is working towards a desired state, cluster can be unavailable.
	ZooKeeperApplyingChangesPhase ZooKeeperOrchestrationPhase = "ApplyingChanges"
	// ZooKeeperMigratingDataPhase ZooKeeper is currently migrating data to another node.
	ZooKeeperMigratingDataPhase ZooKeeperOrchestrationPhase = "MigratingData"
	// ZooKeeperResourceInvalid is marking a resource as invalid, should never happen if admission control is installed correctly.
	ZooKeeperResourceInvalid ZooKeeperOrchestrationPhase = "Invalid"
	ZooKeeperDownScaling     ZooKeeperOrchestrationPhase = "DownScaling"
	ZooKeeperUpScaling       ZooKeeperOrchestrationPhase = "UpScaling"
)

// +kubebuilder:object:root=true
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=zk
// +kubebuilder:categories=zookeeper
// +kubebuilder:printcolumn:name="leader",type="string",JSONPath=".status.leaderNode"
// +kubebuilder:printcolumn:name="nodes",type="integer",JSONPath=".status.availableNodes",description="Available nodes"
// +kubebuilder:printcolumn:name="version",type="string",JSONPath=".spec.version",description="ZooKeeper version"

// Workload is the Schema for the workloads API
type Workload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkloadSpec   `json:"spec,omitempty"`
	Status WorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkloadList contains a list of Workload
type WorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Workload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Workload{}, &WorkloadList{})
}
