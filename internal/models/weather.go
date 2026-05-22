package models

import (
	"time"
)

type WeatherLocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type Weather struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	LocationID    uint      `gorm:"not null;index" json:"location_id"`
	Time          time.Time `json:"time"`
	Temperature   float64   `json:"temperature"`
	Windspeed     float64   `json:"windspeed"`
	Winddirection float64   `json:"winddirection"`
	CreatedAt     time.Time `json:"created_at"`
}

/*
Request:
url := "https://api.open-meteo.com/v1/forecast?latitude=-33.87&longitude=151.21&current_weather=true"
{
	"latitude":-33.848858,
	"longitude":151.19551,
	"generationtime_ms":0.08273124694824219,
	"utc_offset_seconds":0,
	"timezone":"GMT",
	"timezone_abbreviation":"GMT",
	"elevation":69.0,
	"current_weather_units":{
		"time":"iso8601",
		"interval":"seconds",
		"temperature":"°C",
		"windspeed":"km/h",
		"winddirection":"°",
		"is_day":"",
		"weathercode":"wmo code"
	},
	"current_weather":{
		"time":"2026-05-10T00:45",
		"interval":900,
		"temperature":19.0,
		"windspeed":4.3,
		"winddirection":213,
		"is_day":1,
		"weathercode":0
	}
}
*/

type Rainfall struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	LocationID    uint      `gorm:"not null;uniqueIndex:idx_location_time" json:"location_id"`
	Time          time.Time `gorm:"uniqueIndex:idx_location_time" json:"time"`
	Precipitation float64   `json:"precipitation"`
	CreatedAt     time.Time `json:"created_at"`
}

/*
https://api.open-meteo.com/v1/forecast?latitude=-33.8688&longitude=151.2093&start_date=2026-05-10&end_date=2026-05-14&daily=precipitation_sum
{
  "latitude": -33.848858,
  "longitude": 151.19551,
  "generationtime_ms": 0.027060508728027344,
  "utc_offset_seconds": 0,
  "timezone": "GMT",
  "timezone_abbreviation": "GMT",
  "elevation": 86.0,
  "daily_units": {
    "time": "iso8601",
    "precipitation_sum": "mm"
  },
  "daily": {
    "time": [
      "2026-05-10",
      "2026-05-11",
      "2026-05-12",
      "2026-05-13",
      "2026-05-14"
    ],
    "precipitation_sum": [
      0.50,
      3.90,
      3.50,
      11.70,
      6.10
    ]
  }
}
*/
