package repository

import (
	"context"

	"goapi.railway.app/internal/models"
	"gorm.io/gorm"
)

type WeatherRecordRepository interface {
	List(ctx context.Context, id uint) ([]models.WeatherRecord, error)
}

type weatherRecordRepository struct {
	db *gorm.DB
}

func NewWeatherRecordRepository(db *gorm.DB) WeatherRecordRepository {
	return &weatherRecordRepository{db: db}
}

func (r *weatherRecordRepository) List(ctx context.Context, id uint) ([]models.WeatherRecord, error) {
	var WeatherRecords []models.WeatherRecord
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Order("id DESC").
		Find(&WeatherRecords).Error
	return WeatherRecords, err
}

// ========================================
type WeatherRainfallRepository interface {
	List(ctx context.Context) ([]models.Rainfall, error)
}

type weatherRainfallRepository struct {
	db *gorm.DB
}

func NewWeatherRainfallRepository(db *gorm.DB) WeatherRainfallRepository {
	return &weatherRainfallRepository{db: db}
}

func (r *weatherRainfallRepository) List(ctx context.Context) ([]models.Rainfall, error) {
	var Rainfalls []models.Rainfall
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&Rainfalls).Error
	return Rainfalls, err
}
