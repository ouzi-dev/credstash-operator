package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CredstashSecretDef struct {
	Key     string `json:"key,omitempty"`
	Table   string `json:"table,omitempty"`
	Version string `json:"version,omitempty"`
}

// CredstashSecretSpec defines the desired state of CredstashSecret
type CredstashSecretSpec struct {
	SecretName string               `json:"name,omitempty"`
	Secrets    []CredstashSecretDef `json:"secrets,omitempty"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// CredstashSecretStatus defines the observed state of CredstashSecret
type CredstashSecretStatus struct {
	SecretName string `json:"name,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CredstashSecret is the Schema for the credstashsecrets API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=credstashsecrets,scope=Namespaced
type CredstashSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CredstashSecretSpec   `json:"spec,omitempty"`
	Status CredstashSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CredstashSecretList contains a list of CredstashSecret
type CredstashSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CredstashSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CredstashSecret{}, &CredstashSecretList{})
}
