package dto

import "goapi.railway.app/internal/models"

type ReadHistoryRequest struct {
	PlantId uint `json:"plant_id"`
}

type WriteHistoryRequest struct {
	PlantId     uint   `json:"plant_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Agent       string `json:"agent"`
}

type PlantHistoryResponse struct {
	Log []models.PlantHistory `json:"log"`
}
