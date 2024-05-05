package models

import "k8s.io/apimachinery/pkg/runtime"

// Defines a desired runtime component
// +kubebuilder:object:generate=true
type ComponentSpec struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata,omitempty"`
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Properties   runtime.RawExtension `json:"properties,omitempty"`
	Constraints  string               `json:"constraints,omitempty"`
	Dependencies []string             `json:"dependencies,omitempty"`
}
