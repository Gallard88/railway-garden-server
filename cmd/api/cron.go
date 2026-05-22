package main

import (
	"log"
	"time"

	"goapi.railway.app/internal/models"
)

func (app *application) setupCronJobs() {
	// Example 1: Run every minute
	/*
		app.cron.AddFunc("* * * * *", func() {
			log.Println("Running task: Clean up old data")
			app.cleanupOldDataJob()
		})

		// Example 2: Run every day at midnight
		app.cron.AddFunc("0 0 * * *", func() {
			log.Println("Running task: Daily report")
			app.dailyReportJob()
		})

		// Example 3: Run every hour
		app.cron.AddFunc("0 * * * *", func() {
			log.Println("Running task: Update statistics")
			app.updateStatisticsJob()
		})

		// Example 4: Run every Monday at 9 AM
		app.cron.AddFunc("0 9 * * 1", func() {
			log.Println("Running task: Weekly summary")
			app.weeklySummaryJob()
		})
	*/
}

// Example cron job implementations
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

func (app *application) dailyReportJob() {
	task := models.Task{
		Name:        "daily_report",
		Status:      "running",
		Description: "Generate daily statistics",
	}
	app.db.Create(&task)

	task.Status = "success"
	task.Description = "Daily report generated successfully"
	task.LastRunAt = time.Now()
	app.db.Save(&task)

	log.Printf("Daily report: %d active users, %d published posts", 1, 2)
}

func (app *application) updateStatisticsJob() {
	task := models.Task{
		Name:        "update_statistics",
		Status:      "running",
		Description: "Update application statistics",
	}
	app.db.Create(&task)

	// Your statistics update logic here
	// For example: calculate averages, update counters, etc.

	task.Status = "success"
	task.LastRunAt = time.Now()
	app.db.Save(&task)

	log.Println("Statistics updated successfully")
}

func (app *application) weeklySummaryJob() {
	task := models.Task{
		Name:        "weekly_summary",
		Status:      "running",
		Description: "Generate weekly summary",
	}
	app.db.Create(&task)

	// Your weekly summary logic here
	// Could send emails, generate reports, etc.

	task.Status = "success"
	task.LastRunAt = time.Now()
	app.db.Save(&task)

	log.Println("Weekly summary generated")
}
