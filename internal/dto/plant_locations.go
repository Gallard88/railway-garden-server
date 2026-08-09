package dto

import (
	"time"

	"goapi.railway.app/internal/models"
)

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
	ID                    uint                  `json:"id"`
	Name                  string                `json:"name"`
	Zone                  uint                  `json:"zone"`
	WaterFreq             float64               `json:"water_freq"`
	PlantedDate           time.Time             `json:"planted_date"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	LastWatered           time.Time             `json:"last_watered"`
	NextWater             time.Time             `json:"next_water"`
	ContainerType         string                `json:"container_type"`         // in-ground or pot
	SunExposure           string                `json:"sun_exposure"`           // exposed, part shade, full shade
	SoilType              string                `json:"soil_type"`              // sandy, loam, clay
	ET0                   float64               `json:"et0_multiplier"`         // derived from sun exposure (1.0 / 0.6 / 0.3)
	DeficitThreshold      float64               `json:"deficit_threshold"`      // mm of deficit before watering triggered (lower for sandy, higher for clay)
	LookbackDays          uint                  `json:"lookback_days"`          // 1 for pots, 3-5 for in-ground
	RainfallEffectiveness float64               `json:"rainfall_effectiveness"` // how much rainfall actually reaches the zone (pots under eaves etc), default 1.0
	Log                   []models.PlantHistory `json:"log"`
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

type WaterPlanResponse struct {
	ID        uint      `json:"id"`
	NextWater time.Time `json:"planted_date"`
}

// ================================================
type ReadHistoryRequest struct {
	PlantId uint `json:"plant_id"`
}

type WriteHistoryRequest struct {
	PlantId     uint   `json:"plant_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Agent       string `json:"agent"`
}

type PlantHistoryResponse struct {
	Log []models.PlantHistory `json:"log"`
}
