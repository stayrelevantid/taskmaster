package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft Delete
}

// Request body untuk Create & Update Task
type TaskRequest struct {
	Title  string `json:"title" binding:"required,min=3"`
	Status string `json:"status"`
}
