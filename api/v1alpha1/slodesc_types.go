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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TailSLO describes parameters for tail SLO
type TailSLO struct {
	Enable     bool `json:"enable"`
	Percentile uint `json:"percentile,omitempty"`
	Latency    uint `json:"latency"`
}

// MeanSLO describes parameters for mean SLO
type MeanSLO struct {
	Enable bool `json:"enable"`
	From   uint `json:"from,omitempty"`
	To     uint `json:"to,omitempty"`
}

// ShapeSLO describes parameters for shape SLO
type ShapeSLO struct {
	Enable bool `json:"enable"`
	Stddev uint `json:"stddev"`
}

// FaaSTopology describes workflow topology where we ask for the SLOs
type FaaSTopology struct {
	Topology string   `json:"topology,omitempty"`
	Tasks    []string `json:"tasks,omitempty"`
}

// SLODescSpec defines the desired state of SLODesc
type SLODescSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SLODesc. Edit SLODesc_types.go to remove/update
	RateLimit     uint         `json:"ratelimit,omitempty"`
	Tail          TailSLO      `json:"tail"`
	Mean          MeanSLO      `json:"mean"`
	Shape         ShapeSLO     `json:"shape"`
	Workflow      FaaSTopology `json:"workflow,omitempty"`
	SupportPeriod uint         `json:"supportperiod"`
}

type WorkflowStatus struct {
	MeanLatency   uint            `json:"meanlat"`
	TailLatency   uint            `json:"taillat"`
	StddevLatency uint            `json:"stddevlat"`
	InputRate     uint            `json:"inputrate"`
	RunningPods   map[string]uint `json:"runningpods"`
	DesiredPods   map[string]uint `json:"desiredpods"`
}

// SLODescStatus defines the observed state of SLODesc
type SLODescStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Metrics         []WorkflowStatus `json:"metrics"`
	MeanColdStart   uint             `json:"meancoldstart"`
	TailColdStart   uint             `json:"tailcoldstart"`
	StddevColdStart uint             `json:"stddevcoldstart"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SLODesc is the Schema for the slodescs API
type SLODesc struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SLODescSpec   `json:"spec,omitempty"`
	Status SLODescStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SLODescList contains a list of SLODesc
type SLODescList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SLODesc `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SLODesc{}, &SLODescList{})
}
