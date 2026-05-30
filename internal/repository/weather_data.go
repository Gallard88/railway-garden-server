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

const readingsPerHour = 4
const hoursPerDay = 24
const numberOfDays = 7
const maxRecords = readingsPerHour * hoursPerDay * numberOfDays // 48 hours of data

func (r *weatherRecordRepository) List(ctx context.Context, id uint) ([]models.WeatherRecord, error) {
	var WeatherRecords []models.WeatherRecord
	err := r.db.WithContext(ctx).
		Where("location_id = ?", id).
		Order("id DESC").
		Limit(maxRecords).
		Find(&WeatherRecords).Error
	return WeatherRecords, err
}

// ========================================
type WeatherRainfallRepository interface {
	List(ctx context.Context, id uint) ([]models.Rainfall, error)
}

type weatherRainfallRepository struct {
	db *gorm.DB
}

func NewWeatherRainfallRepository(db *gorm.DB) WeatherRainfallRepository {
	return &weatherRainfallRepository{db: db}
}

func (r *weatherRainfallRepository) List(ctx context.Context, id uint) ([]models.Rainfall, error) {
	var Rainfalls []models.Rainfall
	err := r.db.WithContext(ctx).
		Where("location_id = ?", id).
		Order("id DESC").
		Limit(numberOfDays).
		Find(&Rainfalls).Error
	return Rainfalls, err
}
