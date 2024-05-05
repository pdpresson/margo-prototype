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
	DeviceRepo DeviceRepo          `json:"deviceRepo"`
	Sources    []Source            `json:"sources"`
	Properties map[string]Property `json:"properties,omitempty"`
}
