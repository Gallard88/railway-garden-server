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

// Example cron job implementations
/*
func (app *application) cleanupOldDataJob() {
	task := models.Task{
		Name:        "cleanup_old_data",
		Status:      "running",
		Description: "Clean up old deleted records",
	}
	app.db.Create(&task)

	// Delete soft-deleted records older than 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	result := app.db.Unscoped().
		Where("deleted_at < ?", thirtyDaysAgo).
		Delete(&models.User{})

	if result.Error != nil {
		task.Status = "failed"
		task.Description = "Failed: " + result.Error.Error()
	} else {
		task.Status = "success"
		task.Description = "Cleaned up old data successfully"
	}

	task.LastRunAt = time.Now()
	app.db.Save(&task)

	log.Printf("Cleanup job completed. Deleted %d records", result.RowsAffected)
}
*/
