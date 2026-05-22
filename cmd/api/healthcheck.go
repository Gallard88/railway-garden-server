package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/models"
)

func (app *application) healthcheckHandler(c *gin.Context) {
	// Check database connection
	sqlDB, err := app.db.DB()
	dbStatus := "healthy"
	if err != nil {
		dbStatus = "unhealthy"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "unhealthy"
	}

	data := gin.H{
		"status": "available",
		"system_info": gin.H{
			"environment": app.config.env,
			"version":     version,
		},
		"database": gin.H{
			"status": dbStatus,
		},
		"cron": gin.H{
			"active_jobs": len(app.cron.Entries()),
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, data)
}

// ============= TASK HANDLERS (for monitoring cron jobs) =============

func (app *application) listTasksHandler(c *gin.Context) {
	var tasks []models.Task

	result := app.db.Order("created_at DESC").Find(&tasks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": result.RowsAffected,
	})
}

func (app *application) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	result := app.db.First(&task, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}
