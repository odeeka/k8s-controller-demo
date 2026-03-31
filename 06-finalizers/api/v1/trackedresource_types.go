package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// TrackedResourceSpec describes a resource that needs cleanup before deletion.
type TrackedResourceSpec struct {
	// Message is an arbitrary string that the controller will "track".
	// In a real scenario this might be an external resource ID, a lease name, etc.
	// +kubebuilder:validation:Required
	Message string `json:"message"`
}

// TrackedResourceStatus records what the controller has observed.
type TrackedResourceStatus struct {
	// Phase is "Active" when the resource is being tracked normally.
	// "Terminating" when cleanup is in progress.
	Phase string `json:"phase,omitempty"`

	// Message describes the current state.
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".spec.message"

// TrackedResource is the Schema for the trackedresources API.
type TrackedResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrackedResourceSpec   `json:"spec,omitempty"`
	Status TrackedResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TrackedResourceList contains a list of TrackedResource objects.
type TrackedResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrackedResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrackedResource{}, &TrackedResourceList{})
}
