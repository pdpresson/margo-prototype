package models

type Device struct {
	ID       string         `json:"id"`
	Metadata DeviceMetadata `json:"metadata"`
}

type DeviceMetadata struct {
	Name string `json:"name"`
}
