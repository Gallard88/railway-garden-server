package dto

import "time"

// WeatherLocationResponse - What we return to clients
type WeatherLocationResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateWeatherLocationRequest - For creating new locations
type CreateWeatherLocationRequest struct {
	Name      string  `json:"name" binding:"required,min=2"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// ListWeatherLocationsResponse - List all locations
type ListWeatherLocationsResponse struct {
	Locations []WeatherLocationResponse `json:"locations"`
	Count     int                       `json:"count"`
}
