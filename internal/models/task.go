package models

import "time"

// Task model
type Task struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" validate:"required" gorm:"index"`
	Summary     string    `json:"summary" validate:"required"`
	PerformedAt time.Time `json:"performed_at"`
}

// TaskFilter model
type TaskFilter struct {
	Role   string `json:"role"`
	UserID string `json:"user_id"`
}
