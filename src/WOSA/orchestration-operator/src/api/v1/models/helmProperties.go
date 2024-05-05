package models

/*
apiVersion: solution.symphony/v1
kind: Solution
metadata:
  name: hello-world-solution
  namespace: "my-symphony"
spec:
  components:
  - name: hello-world
    type: container
    properties:
      helm.chart.name: "hello-world"
      helm.chart.version: "0.0.1"
      helm.repo: "oci://ghcr.io/pdpresson/charts/hello-world"
      helm.chart.namespace: "my-hello-world"
      helm.values:
        Greeting: "Default"
        Target: "Default"
*/

type HelmPropertyConfig struct {
	Name    string                 `json:"helm.chart.name"`
	Version string                 `json:"helm.chart.version"`
	Repo    string                 `json:"helm.chart.repo"`
	Wait    bool                   `json:"helm.chart.wait"`
	Values  map[string]interface{} `json:"helm.values"`
}
