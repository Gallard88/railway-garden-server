package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/service"
)

type WeatherRecordHandler struct {
	recordService service.WeatherRecordService
}

func NewWeatherRecordHandler(recordService service.WeatherRecordService) *WeatherRecordHandler {
	return &WeatherRecordHandler{
		recordService: recordService,
	}
}

// ListRecords - GET /v1/weather-records
func (h *WeatherRecordHandler) ListRecords(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	// Call service
	result, err := h.recordService.ListRecords(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, result)
}

// RegisterRoutes - Register all weather location routes
func (h *WeatherRecordHandler) RegisterRoutes(router *gin.RouterGroup) {
	locations := router.Group("/weather-records")
	{
		locations.GET("/:id", h.ListRecords) // GET /v1/weather-records
	}
}

// ========================================
type WeatherRainfallHandler struct {
	recordService service.WeatherRainfallService
}

func NewWeatherRainfallHandler(recordService service.WeatherRainfallService) *WeatherRainfallHandler {
	return &WeatherRainfallHandler{
		recordService: recordService,
	}
}

// ListRecords - GET /v1/weather-records
func (h *WeatherRainfallHandler) ListRecords(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	// Call service
	result, err := h.recordService.ListRecords(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, result)
}

// RegisterRoutes - Register all weather location routes
func (h *WeatherRainfallHandler) RegisterRoutes(router *gin.RouterGroup) {
	locations := router.Group("/weather-rainfall")
	{
		locations.GET("/:id", h.ListRecords) // GET /v1/weather-rainfall
	}
}
