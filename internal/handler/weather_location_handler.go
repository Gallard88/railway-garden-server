package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/service"
)

type WeatherLocationHandler struct {
	locationService service.WeatherLocationService
}

func NewWeatherLocationHandler(locationService service.WeatherLocationService) *WeatherLocationHandler {
	return &WeatherLocationHandler{
		locationService: locationService,
	}
}

// GetLocation - GET /v1/weather-locations/:id
func (h *WeatherLocationHandler) GetLocation(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	// 2. Call service
	location, err := h.locationService.GetLocation(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusOK, gin.H{"location": location})
}

// ListLocations - GET /v1/weather-locations
func (h *WeatherLocationHandler) ListLocations(c *gin.Context) {
	// Call service
	result, err := h.locationService.ListLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, result)
}

// CreateLocation - POST /v1/weather-locations
func (h *WeatherLocationHandler) CreateLocation(c *gin.Context) {
	// 1. Parse and validate request
	var req dto.CreateWeatherLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call service
	location, err := h.locationService.CreateLocation(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusCreated, gin.H{"location": location})
}

// RegisterRoutes - Register all weather location routes
func (h *WeatherLocationHandler) RegisterRoutes(router *gin.RouterGroup) {
	locations := router.Group("/weather-locations")
	{
		locations.GET("", h.ListLocations)   // GET /v1/weather-locations
		locations.GET("/:id", h.GetLocation) // GET /v1/weather-locations/:id
		locations.POST("", h.CreateLocation) // POST /v1/weather-locations
	}
}
