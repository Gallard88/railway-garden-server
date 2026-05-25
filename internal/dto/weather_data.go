package dto

import "time"

// =============================================
type WeatherRecordResponse struct {
	ID uint `json:"id"`
	//Name          string    `json:"name"`
	LocationID    uint      `json:"location_id"`
	Time          time.Time `json:"time"`
	Temperature   float64   `json:"temperature"`
	Windspeed     float64   `json:"windspeed"`
	Winddirection float64   `json:"winddirection"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListWeatherLocationsResponse - List all locations
type ListWeatherRecordsResponse struct {
	Records []WeatherRecordResponse `json:"records"`
	Count   int                     `json:"count"`
}

// =============================================
type WeatherRainfallResponse struct {
	ID            uint      `json:"id"`
	LocationID    uint      `json:"location_id"`
	Time          time.Time `json:"time"`
	Precipitation float64   `json:"precipitation"`
	CreatedAt     time.Time `json:"created_at"`
}

// ListWeatherLocationsResponse - List all locations
type ListWeatherRainfallResponse struct {
	Records []WeatherRainfallResponse `json:"records"`
	Count   int                       `json:"count"`
}
