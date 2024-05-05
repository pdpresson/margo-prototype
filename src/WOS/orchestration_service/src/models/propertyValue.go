package models

type PropertyValue struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
