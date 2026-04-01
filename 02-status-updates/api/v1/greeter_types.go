package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// GreeterSpec defines the desired state of a Greeter.
type GreeterSpec struct {
	// Greeting is the word used to greet (e.g. "Hello", "Hi", "Howdy").
	// +kubebuilder:validation:Required
	Greeting string `json:"greeting"`

	// TargetName is who to greet.
	// +kubebuilder:validation:Required
	TargetName string `json:"targetName"`
}

// GreeterStatus defines the observed state of a Greeter.
// The controller writes here after processing the spec.
type GreeterStatus struct {
	// Phase is the current lifecycle state of the Greeter.
	// "Ready" means the controller has processed the spec successfully.
	Phase string `json:"phase,omitempty"`

	// Message is the computed greeting message.
	Message string `json:"message,omitempty"`

	// LastUpdatedTime records when the controller last updated this status.
	LastUpdatedTime metav1.Time `json:"lastUpdatedTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.message"

// Greeter is the Schema for the greeters API.
//
// The +kubebuilder:subresource:status marker tells controller-gen to add
// `subresources: status: {}` to the CRD. This separates the spec and status
// update paths so that a regular Update() cannot overwrite .status.
type Greeter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GreeterSpec   `json:"spec,omitempty"`
	Status GreeterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GreeterList contains a list of Greeter objects.
type GreeterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Greeter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Greeter{}, &GreeterList{})
}
