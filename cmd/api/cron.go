package main

import (
	"context"
	"log"

	"goapi.railway.app/internal/service"
	"gorm.io/gorm"
)

func (app *application) setupCronJobs() {
	ctx := context.Background()

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
	app.fetchRainfallForAllLocations(app.db)

	_, err = app.cron.AddFunc("CRON_TZ=Australia/Sydney 0 6 * * *", func() {
		log.Println("Running task: fetchRainfallForAllLocations() ")
		app.waterService.ProcessRainfallEvent(ctx)
	})
	if err != nil {
		log.Printf("Error scheduling rainfall fetch job: %v", err)
	}
	res, _ := app.waterService.ProcessRainfallEvent(ctx)
	log.Printf("Watering result: %+v", res)

}

func (app *application) fetchWeatherForAllLocations(db *gorm.DB) {
	service.FetchWeatherForAllLocations(db)
}

func (app *application) fetchRainfallForAllLocations(db *gorm.DB) {
	service.FetchRainfallForAllLocations(db)
}
