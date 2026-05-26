package repository

import (
	"context"
	"errors"

	"goapi.railway.app/internal/models"
	"gorm.io/gorm"
)

type PlantZoneRepository interface {
	FindByID(ctx context.Context, id uint) (*models.PlantZone, error)
	List(ctx context.Context) ([]models.PlantZone, error)
	Create(ctx context.Context, zone *models.PlantZone) error
}

type plantZoneRepository struct {
	db *gorm.DB
}

func NewPlantZoneRepository(db *gorm.DB) PlantZoneRepository {
	return &plantZoneRepository{db: db}
}

func (r *plantZoneRepository) FindByID(ctx context.Context, id uint) (*models.PlantZone, error) {
	var zone models.PlantZone
	err := r.db.WithContext(ctx).First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plant zone not found")
		}
		return nil, err
	}
	return &zone, nil
}

func (r *plantZoneRepository) List(ctx context.Context) ([]models.PlantZone, error) {
	var zones []models.PlantZone
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&zones).Error
	return zones, err
}

func (r *plantZoneRepository) Create(ctx context.Context, zone *models.PlantZone) error {
	return r.db.WithContext(ctx).Create(zone).Error
}
