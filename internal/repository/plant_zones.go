package repository

import (
	"context"
	"errors"
	"time"

	"goapi.railway.app/internal/models"
	"gorm.io/gorm"
)

type PlantZoneRepository interface {
	FindByID(ctx context.Context, id uint) (*models.PlantZone, error)
	List(ctx context.Context) ([]models.PlantZone, error)
	ListExposedToRainfall(ctx context.Context) ([]models.PlantZone, error)
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

func (r *plantZoneRepository) ListExposedToRainfall(ctx context.Context) ([]models.PlantZone, error) {
	var zones []models.PlantZone
	err := r.db.WithContext(ctx).
		Where("outdoor = true and rain_threshold > 0").
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
	Water(ctx context.Context, id uint) (*models.Plant, error)
	FindByZone(ctx context.Context, zoneID uint) ([]models.Plant, error)
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

func (r *plantRepository) FindByZone(ctx context.Context, zoneID uint) ([]models.Plant, error) {
	var plants []models.Plant
	err := r.db.WithContext(ctx).
		Where("zone_id = ?", zoneID).
		Find(&plants).Error
	return plants, err
}

func (r *plantRepository) Water(ctx context.Context, id uint) (*models.Plant, error) {
	var plant models.Plant
	err := r.db.WithContext(ctx).First(&plant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plant not found")
		}
		return nil, err
	}
	now := time.Now()
	plant.LastWatered = now
	plant.NextWater = now.AddDate(0, 0, int(plant.WaterFreq))
	plant.UpdatedAt = now

	err = r.db.WithContext(ctx).Save(&plant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plant not found")
		}
		return nil, err
	}

	return &plant, nil
}
