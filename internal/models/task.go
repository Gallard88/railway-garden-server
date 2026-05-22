package models

import "time"

// Task model for cron jobs tracking
type Task struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `gorm:"not null" json:"name"`
	LastRunAt   time.Time `json:"last_run_at"`
	Status      string    `json:"status"` // success, failed, running
	Description string    `json:"description"`
}
