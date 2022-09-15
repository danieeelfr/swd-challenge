package api

import "time"

// Task holds the task values
type Task struct {
	ID          string    `json:"id"`
	Summary     string    `json:"summary" binding:"required"`
	PerformedAt time.Time `json:"performed_at" binding:"required"`
	PerformedBy User      `json:"performed_by"`
}

// TasksResponse holds the tasks list payload
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// User holds the logged user information
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	User     string `json:"user"`
}
