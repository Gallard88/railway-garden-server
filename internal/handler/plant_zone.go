package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/service"
)

type PlantZoneHandler struct {
	zoneService service.PlantZoneService
}

func NewPlantZoneHandler(zoneService service.PlantZoneService) *PlantZoneHandler {
	return &PlantZoneHandler{
		zoneService: zoneService,
	}
}

// GetZone - GET /v1/plant-zones/:id
func (h *PlantZoneHandler) GetZone(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone ID"})
		return
	}

	// 2. Call service
	zone, err := h.zoneService.GetZone(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusOK, gin.H{"zone": zone})
}

// ListZones - GET /v1/plant-zones
func (h *PlantZoneHandler) ListZones(c *gin.Context) {
	// Call service
	result, err := h.zoneService.ListZones(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, result)
}

// CreateZone - POST /v1/plant-zones
func (h *PlantZoneHandler) CreateZone(c *gin.Context) {
	// 1. Parse and validate request
	var req dto.CreatePlantZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call service
	zone, err := h.zoneService.CreateZone(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusCreated, gin.H{"zone": zone})
}

// RegisterRoutes - Register all plant zone routes
func (h *PlantZoneHandler) RegisterRoutes(router *gin.RouterGroup) {
	zones := router.Group("/plant-zones")
	{
		zones.GET("", h.ListZones)   // GET /v1/plant-zones
		zones.GET("/:id", h.GetZone) // GET /v1/plant-zones/:id
		zones.POST("", h.CreateZone) // POST /v1/plant-zones
	}
}
