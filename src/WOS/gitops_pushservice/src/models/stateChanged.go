package models

import "github.com/google/uuid"

type Kind int

const (
	Install Kind = iota + 1
	Update
	Remove
)

type StateChanged struct {
	Kind       Kind                `json:"kind"`
	DeviceId   uuid.UUID           `json:"deviceId"`
	AppId      uuid.UUID           `json:"appId"`
	AppName    string              `json:"appName"`
	DeviceRepo Repo                `json:"deviceRepo"`
	Sources    []Source            `json:"sources"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Repo struct {
	Url    string `json:"url"`
	Branch string `json:"branch"`
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

type Property struct {
	Value   interface{}      `yaml:"value,omitempty"`
	Targets []PropertyTarget `yaml:"targets,omitempty"`
}

type PropertyTarget struct {
	Pointer string `yaml:"pointer"`
}
