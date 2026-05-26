package main

import (
	"log"

	"goapi.railway.app/internal/service"
	"gorm.io/gorm"
)

func (app *application) setupCronJobs() {
	// Example 1: Run every minute

	app.cron.AddFunc("7/15 * * * *", func() {
		//log.Println("Running task: fetchWeatherForAllLocations() ")
		app.fetchWeatherForAllLocations(app.db)
	})
	app.fetchWeatherForAllLocations(app.db)

	_, err := app.cron.AddFunc("CRON_TZ=Australia/Sydney 17 5 * * *", func() {
		log.Println("Running task: fetchRainfallForAllLocations() ")
		app.fetchRainfallForAllLocations(app.db)
	})
	if err != nil {
		log.Printf("Error scheduling rainfall fetch job: %v", err)
	}
	go app.fetchRainfallForAllLocations(app.db)

}

func (app *application) fetchWeatherForAllLocations(db *gorm.DB) {
	service.FetchWeatherForAllLocations(db)
}

func (app *application) fetchRainfallForAllLocations(db *gorm.DB) {
	service.FetchRainfallForAllLocations(db)
}
