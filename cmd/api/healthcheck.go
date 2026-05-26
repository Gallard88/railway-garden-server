package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
