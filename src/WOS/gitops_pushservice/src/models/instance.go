package models

/*
apiVersion: solution.symphony/v1
kind: Instance
metadata:
  name: hello-world-instance
  namespace: my-symphony
spec:
  name: hello-world
  scope: my-symphony
  solution: hello-world-solution
  target:
    name: uknown-target
*/

type Instance struct {
	ApiVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   InstanceMetadata `yaml:"metadata,omitempty"`
	Spec       InstanceSpec     `yaml:"spec,omitempty"`
}

type InstanceMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type InstanceSpec struct {
	Name        string                       `yaml:"name"`
	DisplayName string                       `yaml:"displayName,omitempty"`
	Scope       string                       `yaml:"scope,omitempty"`
	Parameters  map[string]string            `yaml:"parameters,omitempty"`
	Metadata    map[string]string            `yaml:"metadata,omitempty"`
	Solution    string                       `yaml:"solution"`
	Target      TargetSelector               `yaml:"target,omitempty"`
	Arguments   map[string]map[string]string `yaml:"arguments,omitempty"`
	Generation  string                       `yaml:"generation,omitempty"`
	Version     string                       `yaml:"version,omitempty"`
}

type TargetSelector struct {
	Name     string            `yaml:"name,omitempty"`
	Selector map[string]string `yaml:"selector,omitempty"`
}
