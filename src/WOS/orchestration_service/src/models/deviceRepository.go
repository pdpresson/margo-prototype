package models

import "github.com/google/uuid"

type DeviceRepo struct {
	ID     uuid.UUID `json:"id,omitempty"`
	Url    string    `json:"url"`
	Branch string    `json:"branch"`
}
