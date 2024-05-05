package models

// Defines a desired runtime component
type ComponentSpec struct {
	Name         string            `yaml:"name"`
	Type         string            `yaml:"type"`
	Metadata     map[string]string `yaml:"metadata,omitempty"`
	Properties   interface{}       `yaml:"properties,omitempty"`
	Constraints  string            `yaml:"constraints,omitempty"`
	Dependencies []string          `yaml:"dependencies,omitempty"`
}
