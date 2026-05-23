package main

import (
	"log"

	"goapi.railway.app/internal/services"
	"gorm.io/gorm"
)

func (app *application) setupCronJobs() {
	// Example 1: Run every minute

	app.cron.AddFunc("*/15 * * * *", func() {
		log.Println("Running task: Feathcing ")
		app.fetchWeatherForAllLocations(app.db)
	})
	app.fetchWeatherForAllLocations(app.db)

	app.cron.AddFunc("0 20 * * *", func() {
		log.Println("Running task: Feathcing ")
		app.fetchRainfallForAllLocations(app.db)
	})
	app.fetchRainfallForAllLocations(app.db)

}

func (app *application) fetchWeatherForAllLocations(db *gorm.DB) {
	services.FetchWeatherForAllLocations(db)
}

func (app *application) fetchRainfallForAllLocations(db *gorm.DB) {
	services.FetchRainfallForAllLocations(db)
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
