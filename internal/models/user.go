package models

// User model
type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	User     string `json:"user"`
}

// UserFilter filter model
type UserFilter struct {
	Password string `json:"password"`
	User     string `json:"user"`
}
