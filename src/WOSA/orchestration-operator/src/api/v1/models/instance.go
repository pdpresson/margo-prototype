package models

// InstanceSpec defines the spec property of the InstanceState
// +kubebuilder:object:generate=true
type InstanceSpec struct {
	Name        string                       `json:"name"`
	DisplayName string                       `json:"displayName,omitempty"`
	Scope       string                       `json:"scope,omitempty"`
	Parameters  map[string]string            `json:"parameters,omitempty"` //TODO: Do we still need this?
	Metadata    map[string]string            `json:"metadata,omitempty"`
	Solution    string                       `json:"solution"`
	Target      TargetSelector               `json:"target,omitempty"`
	Arguments   map[string]map[string]string `json:"arguments,omitempty"`
	Generation  string                       `json:"generation,omitempty"`
	// Defines the version of a particular resource
	Version string `json:"version,omitempty"`
}

// TargertRefSpec defines the target the instance will deploy to
// +kubebuilder:object:generate=true
type TargetSelector struct {
	Name     string            `json:"name,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}
