package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ConfigSourceSpec defines the desired state of a ConfigSource.
type ConfigSourceSpec struct {
	// Data is the key-value pairs to store in the managed ConfigMap.
	// All values must be strings (this is a Kubernetes requirement for ConfigMaps).
	Data map[string]string `json:"data"`

	// ConfigMapName is the name to give the generated ConfigMap.
	// If not set, the ConfigMap will have the same name as the ConfigSource.
	// +optional
	ConfigMapName string `json:"configMapName,omitempty"`
}

// ConfigSourceStatus defines the observed state of a ConfigSource.
type ConfigSourceStatus struct {
	// Phase is "Ready" once the ConfigMap has been created/updated.
	Phase string `json:"phase,omitempty"`

	// ManagedConfigMap is the name of the ConfigMap being managed.
	ManagedConfigMap string `json:"managedConfigMap,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ConfigMap",type="string",JSONPath=".status.managedConfigMap"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"

// ConfigSource is the Schema for the configsources API.
type ConfigSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSourceSpec   `json:"spec,omitempty"`
	Status ConfigSourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConfigSourceList contains a list of ConfigSource objects.
type ConfigSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigSource{}, &ConfigSourceList{})
}
