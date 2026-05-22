package handlers

import (
	"net/http"
	"strconv"

	"bitbucket.org/ThomasBurns/gardenserver/database"
	"bitbucket.org/ThomasBurns/gardenserver/models"

	"github.com/gin-gonic/gin"
)

// GetPlants returns all plants
func GetPlants(c *gin.Context) {
	var plants []models.Plant

	if err := database.DB.Order("id ASC").Find(&plants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plants"})
		return
	}

	c.JSON(http.StatusOK, plants)
}

// GetPlant returns a single plant by ID
func GetPlant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}

	var plant models.Plant
	if err := database.DB.First(&plant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plant not found"})
		return
	}

	c.JSON(http.StatusOK, plant)
}

// CreatePlant creates a new plant
func CreatePlant(c *gin.Context) {
	var req models.CreatePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plant := models.Plant{
		Name:      req.Name,
		Zone:      req.Zone,
		WaterFreq: req.WaterFreq,
	}

	if err := database.DB.Create(&plant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plant"})
		return
	}

	c.JSON(http.StatusCreated, plant)
}

// UpdatePlant updates an existing plant
func UpdatePlant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}

	var plant models.Plant
	if err := database.DB.First(&plant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
		return
	}

	var req models.UpdatePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if err := database.DB.Model(&plant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plant"})
		return
	}

	// Reload to get updated timestamps
	database.DB.First(&plant.ID, id)

	c.JSON(http.StatusOK, plant)
}

// DeletePlant deletes a plant
func DeletePlant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}

	result := database.DB.Delete(&models.Plant{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plant"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "plant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plant deleted successfully"})
}
