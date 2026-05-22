package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"goapi.railway.app/internal/database"
	"goapi.railway.app/internal/models"
	"gorm.io/gorm/clause"
)

type OpenMeteoTime struct {
	time.Time
}

func (t *OpenMeteoTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // Remove quotes

	// Open-Meteo uses "2006-01-02T15:04" format (no seconds, no timezone)
	parsed, err := time.Parse("2006-01-02T15:04", s)
	if err != nil {
		return err
	}

	t.Time = parsed
	return nil
}

type OpenMeteoWeatherResponse struct {
	CurrentWeather struct {
		Time          OpenMeteoTime `json:"time"`
		Temperature   float64       `json:"temperature"`
		Windspeed     float64       `json:"windspeed"`
		Winddirection float64       `json:"winddirection"`
	} `json:"current_weather"`
}

// FetchWeatherForAllLocations polls weather for all locations
func FetchWeatherForAllLocations() {
	log.Println("Starting weather fetch for all locations...")

	var locations []models.WeatherLocation
	if err := database.DB.Find(&locations).Error; err != nil {
		log.Printf("Error fetching locations: %v", err)
		return
	}
	for _, location := range locations {
		if err := fetchAndStoreWeather(location); err != nil {
			log.Printf("Error fetching weather for location %d: %v", location.ID, err)
			continue
		}
		log.Printf("Successfully fetched weather for %s", location.Name)
	}

	log.Println("Weather fetch completed")
}

func fetchAndStoreWeather(location models.WeatherLocation) error {
	// Call Open-Meteo API
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
		location.Latitude,
		location.Longitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var weatherData OpenMeteoWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Store in database
	weather := models.Weather{
		LocationID:    location.ID,
		Time:          weatherData.CurrentWeather.Time.Time,
		Temperature:   weatherData.CurrentWeather.Temperature,
		Windspeed:     weatherData.CurrentWeather.Windspeed,
		Winddirection: weatherData.CurrentWeather.Winddirection,
		CreatedAt:     time.Now(),
	}

	if err := database.DB.Create(&weather).Error; err != nil {
		return fmt.Errorf("failed to save weather: %w", err)
	}

	return nil
}

// Custom type for date-only parsing
type Date time.Time

// UnmarshalJSON for automatic parsing during json.Unmarshal
func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	s = s[1 : len(s)-1]

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// Convert back to time.Time for use
func (d Date) Time() time.Time {
	return time.Time(d)
}

type OpenMeteoRainfallResponse struct {
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	GenerationTimeMs float64    `json:"generationtime_ms"`
	UtcOffsetSeconds int        `json:"utc_offset_seconds"`
	Timezone         string     `json:"timezone"`
	TimezoneAbbr     string     `json:"timezone_abbreviation"`
	Elevation        float64    `json:"elevation"`
	DailyUnits       DailyUnits `json:"daily_units"`
	Daily            DailyData  `json:"daily"`
}

type DailyUnits struct {
	Time             string `json:"time"`
	PrecipitationSum string `json:"precipitation_sum"`
}

type DailyData struct {
	Time             []Date    `json:"time"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
}

// FetchRainfallForAllLocations polls weather for all locations
func FetchRainfallForAllLocations() {
	log.Println("Starting rain fetch for all locations...")

	var locations []models.WeatherLocation
	if err := database.DB.Find(&locations).Error; err != nil {
		log.Printf("Error fetching locations: %v", err)
		return
	}
	for _, location := range locations {
		if err := fetchAndStoreRainfall(location); err != nil {
			log.Printf("Error fetching rain for location %d: %v", location.ID, err)
			continue
		}
		log.Printf("Successfully rain weather for %s", location.Name)
	}

	log.Println("Rainfall fetch completed")
}

func fetchAndStoreRainfall(location models.WeatherLocation) error {

	// Start date: 3 days ago
	now := time.Now()
	startDate := now.AddDate(0, 0, -2).Format("2006-01-02")
	endDate := now.AddDate(0, 0, -1).Format("2006-01-02")
	// Call Open-Meteo API
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&start_date=%s&end_date=%s&daily=precipitation_sum",
		location.Latitude,
		location.Longitude,
		startDate,
		endDate,
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response OpenMeteoRainfallResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	fmt.Printf("Rainfall: %+v\n", response)
	// Now you can use parsedTimes
	for i, t := range response.Daily.Time {
		rainfall := response.Daily.PrecipitationSum[i]
		fmt.Printf("%s: %.1f mm\n", t.Time().Format("2006-01-02"), rainfall)

		r := models.Rainfall{
			LocationID:    location.ID,
			Time:          t.Time(),
			Precipitation: rainfall,
			CreatedAt:     time.Now(),
		}
		if err := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "location_id"}, {Name: "time"}},
			DoUpdates: clause.AssignmentColumns([]string{"precipitation"}),
		}).Create(&r).Error; err != nil {
			return fmt.Errorf("failed to save rainfall: %w", err)
		}
	}
	return nil
}
