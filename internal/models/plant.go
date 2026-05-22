package models

import "time"

type Plant struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Zone        uint      `gorm:"check:zone >= 0" json:"zone"`
	WaterFreq   float64   `gorm:"check:water_freq > 0" json:"water_freq"`
	PlantedDate time.Time `json:"planted_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastWatered time.Time `json:"last_watered"`
	NextWater   time.Time `json:"next_water"`
}

type CreatePlantRequest struct {
	Name      string  `json:"name" binding:"required"`
	Zone      uint    `json:"zone" binding:"required,zone"`
	WaterFreq float64 `json:"water_freq" binding:"water_freq"`
}

type UpdatePlantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
}
