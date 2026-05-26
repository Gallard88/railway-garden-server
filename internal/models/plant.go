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

type PlantZone struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	LocationID    uint      `gorm:"check:location_id >= 0" json:"location_id"`
	Outdoor       bool      `json:"outdoor"`
	RainThreshold float64   `json:"rain_threshold"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
