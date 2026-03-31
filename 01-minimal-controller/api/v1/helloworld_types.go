package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// HelloWorldSpec defines the *desired* state of a HelloWorld resource.
// This is what the user writes in the YAML under `spec:`.
type HelloWorldSpec struct {
	// Name is the name to greet. Required.
	Name string `json:"name"`
}

// HelloWorldStatus defines the *observed* state of a HelloWorld resource.
// The controller writes here to report what it has actually done.
// In this first step, the status is intentionally empty — we introduce it in step 02.
type HelloWorldStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced

// HelloWorld is a custom resource that represents a thing to greet.
// The lines above (starting with `// +kubebuilder:`) are *markers* used by
// controller-gen to generate CRD YAML and DeepCopy code. In this repo we
// write those files by hand to keep things transparent.
type HelloWorld struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelloWorldSpec   `json:"spec,omitempty"`
	Status HelloWorldStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HelloWorldList is a list of HelloWorld objects.
// controller-runtime needs this type so it can list resources.
type HelloWorldList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelloWorld `json:"items"`
}

// init() registers our types with SchemeBuilder when the package is loaded.
// This is called automatically; you don't invoke it yourself.
func init() {
	SchemeBuilder.Register(&HelloWorld{}, &HelloWorldList{})
}
