package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// AppDeploymentSpec describes the application to deploy.
type AppDeploymentSpec struct {
	// Image is the container image to run (e.g. "nginx:1.25").
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Replicas is the desired number of pod replicas.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// Port is the container port to expose.
	// +optional
	// +kubebuilder:default=80
	Port int32 `json:"port,omitempty"`
}

// AppDeploymentStatus records what the controller has observed.
type AppDeploymentStatus struct {
	// Phase is "Available" when all desired replicas are ready,
	// "Progressing" otherwise.
	Phase string `json:"phase,omitempty"`

	// AvailableReplicas mirrors the managed Deployment's availableReplicas.
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Available",type="integer",JSONPath=".status.availableReplicas"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"

// AppDeployment is the Schema for the appdeployments API.
type AppDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppDeploymentSpec   `json:"spec,omitempty"`
	Status AppDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppDeploymentList contains a list of AppDeployment objects.
type AppDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppDeployment{}, &AppDeploymentList{})
}
