package repository

import (
	"context"
	"errors"

	"goapi.railway.app/internal/models"
	"gorm.io/gorm"
)

type WeatherLocationRepository interface {
	FindByID(ctx context.Context, id uint) (*models.WeatherLocation, error)
	List(ctx context.Context) ([]models.WeatherLocation, error)
	Create(ctx context.Context, location *models.WeatherLocation) error
}

type weatherLocationRepository struct {
	db *gorm.DB
}

func NewWeatherLocationRepository(db *gorm.DB) WeatherLocationRepository {
	return &weatherLocationRepository{db: db}
}

func (r *weatherLocationRepository) FindByID(ctx context.Context, id uint) (*models.WeatherLocation, error) {
	var location models.WeatherLocation
	err := r.db.WithContext(ctx).First(&location, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("location not found")
		}
		return nil, err
	}
	return &location, nil
}

func (r *weatherLocationRepository) List(ctx context.Context) ([]models.WeatherLocation, error) {
	var locations []models.WeatherLocation
	err := r.db.WithContext(ctx).
		Order("name DESC").
		Find(&locations).Error
	return locations, err
}

func (r *weatherLocationRepository) Create(ctx context.Context, location *models.WeatherLocation) error {
	return r.db.WithContext(ctx).Create(location).Error
}
