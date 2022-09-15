package api

// LoginResponse holds the login response payload
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
