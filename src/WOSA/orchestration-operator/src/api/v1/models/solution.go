package models

// +kubebuilder:object:generate=true
type SolutionSpec struct {
	DisplayName string            `json:"displayName,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Components  []ComponentSpec   `json:"components,omitempty"`
	// Defines the version of a particular resource
	Version string `json:"version,omitempty"`
}
