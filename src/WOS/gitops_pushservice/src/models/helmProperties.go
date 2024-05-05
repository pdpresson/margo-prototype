package models

type HelmPropertyConfig struct {
	Name    string                 `yaml:"helm.chart.name"`
	Version string                 `yaml:"helm.chart.version"`
	Repo    string                 `yaml:"helm.chart.repo"`
	Wait    bool                   `yaml:"helm.chart.wait"`
	Values  map[string]interface{} `yaml:"helm.values"`
}
