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

// ================================================
type PlantRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Plant, error)
	List(ctx context.Context) ([]models.Plant, error)
	Create(ctx context.Context, plant *models.Plant) error
	DeletePlant(ctx context.Context, id uint) error
}

type plantRepository struct {
	db *gorm.DB
}

func NewPlantRepository(db *gorm.DB) PlantRepository {
	return &plantRepository{db: db}
}

func (r *plantRepository) FindByID(ctx context.Context, id uint) (*models.Plant, error) {
	var plant models.Plant
	err := r.db.WithContext(ctx).First(&plant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plant not found")
		}
		return nil, err
	}
	return &plant, nil
}

func (r *plantRepository) List(ctx context.Context) ([]models.Plant, error) {
	var plants []models.Plant
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&plants).Error
	return plants, err
}

func (r *plantRepository) Create(ctx context.Context, plant *models.Plant) error {
	return r.db.WithContext(ctx).Create(plant).Error
}

func (r *plantRepository) DeletePlant(ctx context.Context, id uint) error {
	var plant models.Plant
	err := r.db.WithContext(ctx).First(&plant, id).Error
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&plant).Error
}
