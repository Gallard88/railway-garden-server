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
	LocationID    uint    `json:"location_id"`
	Outdoor       bool    `json:"outdoor"`
	RainThreshold float64 `json:"rain_threshold"`
}

// ListPlantZonesResponse - List all plant zones
type ListPlantZonesResponse struct {
	Zones []PlantZoneResponse `json:"zones"`
	Count int                 `json:"count"`
}

// ================================================
type PlantResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Zone        uint      `json:"zone"`
	WaterFreq   float64   `json:"water_freq"`
	PlantedDate time.Time `json:"planted_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastWatered time.Time `json:"last_watered"`
	NextWater   time.Time `json:"next_water"`
}

// CreatePlantRequest - For creating new plants
type CreatePlantRequest struct {
	Name        string    `json:"name" binding:"required,min=2"`
	Zone        uint      `json:"zone"`
	WaterFreq   float64   `json:"water_freq"`
	PlantedDate time.Time `json:"planted_date"`
}

// ListPlantsResponse - List all plants
type ListPlantsResponse struct {
	Plants []PlantResponse `json:"plants"`
	Count  int             `json:"count"`
}
