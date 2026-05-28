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

// ===============================================
type PlantHandler struct {
	plantService service.PlantService
}

func NewPlantHandler(plantService service.PlantService) *PlantHandler {
	return &PlantHandler{
		plantService: plantService,
	}
}

// GetPlant - GET /v1/plants/:id
func (h *PlantHandler) GetPlant(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant ID"})
		return
	}

	// 2. Call service
	plant, err := h.plantService.GetPlant(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusOK, gin.H{"plant": plant})
}

// ListPlants - GET /v1/plants
func (h *PlantHandler) ListPlants(c *gin.Context) {
	// Call service
	result, err := h.plantService.ListPlants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusOK, result)
}

// CreatePlant - POST /v1/plants
func (h *PlantHandler) CreatePlant(c *gin.Context) {
	// 1. Parse and validate request
	var req dto.CreatePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call service
	plant, err := h.plantService.CreatePlant(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusCreated, gin.H{"plant": plant})
}

// GetPlant - GET /v1/plants/:id
func (h *PlantHandler) DeletePlant(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant ID"})
		return
	}

	// 2. Call service
	plant, err := h.plantService.GetPlant(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	err = h.plantService.DeletePlant(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plant deleted successfully", "plant": plant})
}

// RegisterRoutes - Register all plant zone routes
func (h *PlantHandler) RegisterRoutes(router *gin.RouterGroup) {
	zones := router.Group("/plants")
	{
		zones.GET("", h.ListPlants)         // GET /v1/plants
		zones.GET("/:id", h.GetPlant)       // GET /v1/plants/:id
		zones.POST("", h.CreatePlant)       // POST /v1/plants
		zones.DELETE("/:id", h.DeletePlant) // DELETE /v1/plants/:id
		//		zones.PUT("/:id/", h.UpdatePlant)     // PUT /v1/plants/:id/mark-watered
		//		zones.PUT("/:id/water", h.WaterPlant) // PUT /v1/plants/:id/mark-watered
	}
}
