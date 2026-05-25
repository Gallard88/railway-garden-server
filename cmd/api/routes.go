package main

import (
	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/handler"
	"goapi.railway.app/internal/repository"
	"goapi.railway.app/internal/service"
)

func (app *application) routes() *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware (customize as needed)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize repositories
	locationRepo := repository.NewWeatherLocationRepository(app.db)

	// Initialize services
	locationService := service.NewWeatherLocationService(locationRepo)

	// Initialize handlers
	locationHandler := handler.NewWeatherLocationHandler(locationService)

	// Health check route
	router.GET("/v1/healthcheck", app.healthcheckHandler)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Cron/Task routes (for monitoring scheduled jobs)
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", app.listTasksHandler)
			tasks.GET("/:id", app.getTaskHandler)
		}
	}
	locationHandler.RegisterRoutes(v1)

	return router
}
