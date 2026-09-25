package dto

import "time"

// CreateWeatherLocationRequest - For creating new locations
type CreateWeatherLocationRequest struct {
	Name      string  `json:"name" binding:"required,min=2"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// ListWeatherLocationsResponse - List all locations
type ListLocationResponse struct {
	Locations []LocationResponse `json:"locations"`
	Count     int                `json:"count"`
}

// =============================================
type LocationResponse struct {
	ID          uint                      `gorm:"primaryKey" json:"id"`
	Name        string                    `gorm:"not null" json:"name"`
	Latitude    float64                   `json:"latitude"`
	Longitude   float64                   `json:"longitude"`
	Temperature []WeatherRecordResponse   `json:"temperature"`
	Rainfall    []WeatherRainfallResponse `json:"rainfall"`
}

// =============================================
type WeatherRecordResponse struct {
	ID            uint      `json:"id"`
	Time          time.Time `json:"time"`
	Temperature   float64   `json:"temperature"`
	Windspeed     float64   `json:"windspeed"`
	Winddirection float64   `json:"winddirection"`
}

// =============================================
type WeatherRainfallResponse struct {
	ID                 uint      `json:"id"`
	Time               time.Time `json:"time"`
	Precipitation      float64   `json:"precipitation"`
	Evapotranspiration float64   `json:"evapotranspiration"`
}

// =============================================
// =============================================
// ListWeatherLocationsResponse - List all locations
type ListWeatherRecordsResponse struct {
	Records []WeatherRecordResponse `json:"records"`
	Count   int                     `json:"count"`
}

// ListWeatherLocationsResponse - List all locations
type ListWeatherRainfallResponse struct {
	Records []WeatherRainfallResponse `json:"records"`
	Count   int                       `json:"count"`
}
