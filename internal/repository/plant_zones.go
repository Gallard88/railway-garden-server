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
	FindByZoneID(ctx context.Context, zoneId uint) ([]models.Plant, error)
	List(ctx context.Context) ([]models.Plant, error)
	Create(ctx context.Context, plant *models.Plant) error
	Update(ctx context.Context, updatedPlant models.Plant) error
	DeletePlant(ctx context.Context, id uint) error
	Water(ctx context.Context, id uint) (*models.Plant, error)
	FindByZone(ctx context.Context, zoneID uint) ([]models.Plant, error)
	MarkNeedsWatering(ctx context.Context, plantId uint, reason, description string) error

	// History related functions
	ReadHistory(ctx context.Context, plantId uint) ([]models.PlantHistory, error)
	CreateHistoryEvent(ctx context.Context, event *models.PlantHistory) error
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

func (r *plantRepository) FindByZoneID(ctx context.Context, zoneId uint) ([]models.Plant, error) {
	var plants []models.Plant
	err := r.db.WithContext(ctx).
		Where("zone = ?", zoneId).
		Find(&plants).Error
	return plants, err
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

func (r *plantRepository) MarkNeedsWatering(ctx context.Context, plantId uint, reason, description string) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&models.Plant{}).
		Where("id = ?", plantId).
		Updates(map[string]interface{}{
			"needs_watering": true,
			"updated_at":     now,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("plant not found")
	}
	return r.db.WithContext(ctx).Create(&models.PlantHistory{
		PlantId:     plantId,
		CreatedAt:   time.Now(),
		Name:        reason,
		Description: description,
		Agent:       "CRON",
	}).Error

}

func (r *plantRepository) Update(ctx context.Context, updatedPlant models.Plant) error {
	err := r.db.WithContext(ctx).Save(&updatedPlant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("plant failed to update")
		}
		return err
	}
	return nil

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
	plant.NeedsWatering = false

	err = r.db.WithContext(ctx).Save(&plant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plant not found")
		}
		return nil, err
	}

	return &plant, nil
}

func (p *plantRepository) ReadHistory(ctx context.Context, plantId uint) ([]models.PlantHistory, error) {
	var history []models.PlantHistory
	err := p.db.WithContext(ctx).
		Where("plant_id = ?", plantId).
		Order("created_at desc").
		Find(&history).Error
	return history, err
}

func (p *plantRepository) CreateHistoryEvent(ctx context.Context, event *models.PlantHistory) error {
	// TODO: Verify that plant id is real!
	return p.db.WithContext(ctx).Create(event).Error
}
