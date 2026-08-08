package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/service"
)

// ===============================================
type PlantHistoryHandler struct {
	historyService service.PlantHistoryService
}

func NewPlantHistoryHandler(historyService service.PlantHistoryService) *PlantHistoryHandler {
	return &PlantHistoryHandler{
		historyService: historyService,
	}
}

// GetPlant - GET /v1/plants/:id
func (h *PlantHistoryHandler) GetHistory(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant ID"})
		return
	}

	// 2. Call service
	history, err := h.historyService.ReadHistory(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusOK, gin.H{"plant": id, "log": history.Log})
}

func (h *PlantHistoryHandler) CreateHistory(c *gin.Context) {
	// 1. Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plant ID"})
		return
	}

	var history dto.WriteHistoryRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to parse response: %w", err)})
		return
	}
	history.PlantId = uint(id)

	fmt.Printf("jsonData: %+v\n", c.Request.Body)
	fmt.Printf("history: %+v\n", history)

	// 2. Call service
	updatedPlant, err := h.historyService.CreateHistoryEvent(c.Request.Context(), history)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 3. Return response
	c.JSON(http.StatusOK, gin.H{"plant": updatedPlant})
}

// RegisterRoutes - Register all plant zone routes
func (h *PlantHistoryHandler) RegisterRoutes(router *gin.RouterGroup) {
	zones := router.Group("/plants")
	{
		zones.GET("/:id/history", h.GetHistory)     // GET  /v1/plants/:id/history
		zones.POST("/:id/history", h.CreateHistory) // POST /v1/plants/:id/history
	}
}
