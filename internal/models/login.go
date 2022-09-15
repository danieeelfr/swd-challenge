package models

// Login model
type Login struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
}
