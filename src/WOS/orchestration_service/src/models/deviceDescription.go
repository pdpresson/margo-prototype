package models

import "github.com/google/uuid"

type Device struct {
	Description DeviceDescription
	Apps        []uuid.UUID
}

type DeviceDescription struct {
	ID         uuid.UUID      `json:"id"`
	RepoId     uuid.UUID      `json:"repoId"`
	ApiVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Metadata   DeviceMetadata `json:"metadata"`
	Spec       DeviceSpec     `json:"spec"`
}

type DeviceMetadata struct {
	Annotations DeviceAnnotations `json:"annotations"`
	Name        string            `json:"name"`
}

type DeviceAnnotations struct {
	Platform        string `json:"margo.org/platform"`
	PlatformVersion string `json:"margo.org/platform-version"`
	Vendor          string `json:"margo.org/vendor"`
	Model           string `json:"margo.org/model"`
	SerialNumber    string `json:"margo.org/serial-number"`
	CPUArchitecture string `json:"margo.org/cpu-architecture"`
	VirtualCPUs     int16  `json:"margo.org/vcpus"`
	Memory          string `json:"margo.org/memory"`
	StorageCapacity string `json:"margo.org/storage-capacity"`
}

type DeviceSpec struct {
	NetworkInterfaces []NetworkInterface `json:"network-interfaces"`
	Periferals        Periferal          `json:"periferals"`
}

type NetworkInterface struct {
	Kind string `json:"kind"`
}

type Periferal struct {
	GPUs []GPU `json:"gpus"`
}

type GPU struct {
	Kind   string `json:"kind"`
	Brand  string `json:"brand"`
	Model  string `json:"model"`
	Memory string `json:"memory"`
}
