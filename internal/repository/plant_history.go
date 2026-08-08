package repository

import (
	"context"

	"goapi.railway.app/internal/models"
	"gorm.io/gorm"
)

//type PlantHistory struct {

type PlantHistoryRepository interface {
	ReadHistory(ctx context.Context, plantId uint) ([]models.PlantHistory, error)
	CreateHistoryEvent(ctx context.Context, event *models.PlantHistory) error
}

type plantHistoryRepository struct {
	db *gorm.DB
}

func NewPlantHistoryRepository(db *gorm.DB) PlantHistoryRepository {
	return &plantHistoryRepository{db: db}
}

func (p *plantHistoryRepository) ReadHistory(ctx context.Context, plantId uint) ([]models.PlantHistory, error) {
	var history []models.PlantHistory
	err := p.db.WithContext(ctx).
		Where("plant_id = ?", plantId).
		Order("created_at desc").
		Find(&history).Error
	return history, err
}

func (p *plantHistoryRepository) CreateHistoryEvent(ctx context.Context, event *models.PlantHistory) error {
	// TODO: Verify that plant id is real!
	return p.db.WithContext(ctx).Create(event).Error
}
