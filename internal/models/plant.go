package models

import "time"

type Plant struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	Name                  string    `gorm:"not null" json:"name"`
	Zone                  uint      `gorm:"check:zone >= 0" json:"zone"`
	WaterFreq             float64   `gorm:"check:water_freq > 0" json:"water_freq"`
	PlantedDate           time.Time `json:"planted_date"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	LastWatered           time.Time `json:"last_watered"`
	NextWater             time.Time `json:"next_water"`
	ContainerType         string    `json:"container_type"`         // in-ground or pot
	SunExposure           string    `json:"sun_exposure"`           // exposed, part shade, full shade
	SoilType              string    `json:"soil_type"`              // sandy, loam, clay
	ET0                   float64   `json:"et0_multiplier"`         // derived from sun exposure (1.0 / 0.6 / 0.3)
	DeficitThreshold      float64   `json:"deficit_threshold"`      // mm of deficit before watering triggered (lower for sandy, higher for clay)
	LookbackDays          uint      `json:"lookback_days"`          // 1 for pots, 3-5 for in-ground
	RainfallEffectiveness float64   `json:"rainfall_effectiveness"` // how much rainfall actually reaches the zone (pots under eaves etc), default 1.0
	OverheatedTemp        float64   `json:"overheated_temp"`        // A threshold for determining when a platn has gotten too hot and should be water at the earliest convenience.
	NeedsWatering         bool      `json:"needs_watering"`
}

const (
	ContainerType_Pot    = "pot"
	ContainerType_Ground = "ground"

	Sun_FullSun   = "full-sun"
	Sun_PartSun   = "part-sun"
	Sun_FullShade = "full-shade"

	Soil_Sandy = "sandy"
	Soil_Clay  = "clay"
	Soil_Loam  = "loam"
)

type PlantZone struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	LocationID     uint      `gorm:"check:location_id >= 0" json:"location_id"`
	Outdoor        bool      `json:"outdoor"`
	RainThreshold  float64   `json:"rain_threshold"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	OverheatedTemp float64   `json:"overheated_temp"` // A threshold for determining when a platn has gotten too hot and should be water at the earliest convenience.
}

type PlantHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PlantId     uint      `json:"plant_id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Agent       string    `json:"agent"`
}
