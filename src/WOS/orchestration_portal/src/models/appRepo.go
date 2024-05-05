package models

import "github.com/google/uuid"

type AppRepo struct {
	ID     uuid.UUID `json:"id"`
	Url    string    `json:"url"`
	Branch string    `json:"branch"`
}
