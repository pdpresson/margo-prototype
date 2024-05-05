package models

import "github.com/google/uuid"

type AppDescription struct {
	ID               uuid.UUID           `yaml:"id,omitempty"`
	RepoId           uuid.UUID           `yaml:"repoId,omitempty"`
	ApiVersion       string              `yaml:"apiVersion"`
	Kind             string              `yaml:"kind"`
	Metadata         Metadata            `yaml:"metadata"`
	Sources          []Source            `yaml:"sources"`
	MinimumResources MinimumResources    `yaml:"minimumResourceRequirements"`
	Properties       map[string]Property `yaml:"properties,omitempty"`
	Configuration    Configuration       `yaml:"configuration"`
}

type Metadata struct {
	Id          string  `yaml:"id"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Version     string  `yaml:"version"`
	Catalog     Catalog `yaml:"catalog"`
}

type Catalog struct {
	Application  Application  `yaml:"application"`
	Author       Author       `yaml:"author"`
	Organization Organization `yaml:"organization"`
}

type Application struct {
	Icon         string `yaml:"icon"`
	Tagline      string `yaml:"tagline"`
	Description  string `yaml:"descriptionLog"`
	ReleaseNotes string `yml:"releaseNotes"`
	LicenseFile  string `yaml:"licenseFile"`
	Site         string `yaml:"site"`
}

type Author struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Organization struct {
	Name string `yaml:"name"`
	Site string `yaml:"site"`
}

type Source struct {
	Name       string         `yaml:"name"`
	Type       string         `yaml:"type"`
	Properties SourceProperty `yaml:"properties"`
}

type SourceProperty struct {
	Repository string `yaml:"repository,omitempty"`
	Revision   string `yaml:"revision,omitempty"`
	Wait       bool   `yaml:"wait,omitempty"`
}

type MinimumResources struct {
	CPU     float32          `yaml:"cpu"`
	Memory  float32          `yaml:"memory"`
	Storage StorageResources `yaml:"storage"`
}

type StorageResources struct {
	Containers   float32 `yaml:"containers"`
	Applications float32 `yaml:"appStorage"`
}

type Property struct {
	Value   interface{}      `yaml:"value,omitempty"`
	Targets []PropertyTarget `yaml:"targets,omitempty"`
}

type PropertyTarget struct {
	Pointer string `yaml:"pointer"`
}

type Configuration struct {
	Sections []ConfigSection `yaml:"sections"`
	Schema   []ConfigSchema  `yaml:"schema"`
}

type ConfigSection struct {
	Name     string            `yaml:"name"`
	Settings []SectionSettings `yaml:"settings"`
}

type SectionSettings struct {
	Property    string `yaml:"property"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	InputType   string `yaml:"inputType"`
	Schema      string `yaml:"schema"`
}

type ConfigSchema struct {
	Name       string `yaml:"name"`
	AppliesTo  string `yaml:"appliesTo"`
	MaxLength  int    `yaml:"maxLength,omitempty"`
	AllowEmpty bool   `yaml:"allowEmpty"`
}
