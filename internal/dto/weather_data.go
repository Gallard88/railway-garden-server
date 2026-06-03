package dto

import "time"

// =============================================
type WeatherRecordResponse struct {
	ID            uint      `json:"id"`
	Time          time.Time `json:"time"`
	Temperature   float64   `json:"temperature"`
	Windspeed     float64   `json:"windspeed"`
	Winddirection float64   `json:"winddirection"`
}

// ListWeatherLocationsResponse - List all locations
type ListWeatherRecordsResponse struct {
	Records []WeatherRecordResponse `json:"records"`
	Count   int                     `json:"count"`
}

// =============================================
type WeatherRainfallResponse struct {
	ID                 uint      `json:"id"`
	Time               time.Time `json:"time"`
	Precipitation      float64   `json:"precipitation"`
	Evapotranspiration float64   `json:"evapotranspiration"`
}

// ListWeatherLocationsResponse - List all locations
type ListWeatherRainfallResponse struct {
	Records []WeatherRainfallResponse `json:"records"`
	Count   int                       `json:"count"`
}
