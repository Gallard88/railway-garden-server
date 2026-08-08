package models

import "time"

//History Log
type PlantHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PlantId     uint      `json:"plant_id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Agent       string    `json:"agent"`
}
