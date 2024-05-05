package models

type ResponseError struct {
	Error  string `json:"error"`
	Status int    `json:"-"`
}
