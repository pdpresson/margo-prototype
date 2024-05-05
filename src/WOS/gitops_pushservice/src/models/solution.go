package models

/*
apiVersion: solution.symphony/v1
kind: Solution
metadata:
  name: hello-world-solution
  namespace: my-symphony
spec:
  components:
  - name: hello-world
    type: container
    properties:
      helm.chart.name: "hello-world"
      helm.chart.version: "0.0.1"
      helm.chart.repo: "oci://ghcr.io/pdpresson/charts/hello-world"
      helm.chart.wait: true
      helm.values:
        env:
          APP_GREETING: "Hello"
          APP_TARGET: "World"
*/

type Solution struct {
	ApiVersion string           `yaml:"apiVersion"`
	Kind       string           `yaml:"kind"`
	Metadata   SolutionMetadata `yaml:"metadata,omitempty"`
	Spec       SolutionSpec     `yaml:"spec,omitempty"`
}

type SolutionMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type SolutionSpec struct {
	DisplayName string            `yaml:"displayName,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
	Components  []ComponentSpec   `yaml:"components,omitempty"`
	Version     string            `yaml:"version,omitempty"`
}
