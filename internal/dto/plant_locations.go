package dto

import "time"

// PlantZoneResponse - What we return to clients
type PlantZoneResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	LocationID    uint      `json:"location_id"`
	Outdoor       bool      `json:"outdoor"`
	RainThreshold float64   `json:"rain_threshold"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreatePlantZoneRequest - For creating new plant zones
type CreatePlantZoneRequest struct {
	Name          string  `json:"name" binding:"required,min=2"`
	LocationID    uint    `json:"location_id" binding:"required"`
	Outdoor       bool    `json:"outdoor" binding:"required"`
	RainThreshold float64 `json:"rain_threshold" binding:"required"`
}

// ListPlantZonesResponse - List all plant zones
type ListPlantZonesResponse struct {
	Zones []PlantZoneResponse `json:"zones"`
	Count int                 `json:"count"`
}
