package service

import (
	"context"
	"errors"
	"fmt"

	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/models"
	"goapi.railway.app/internal/repository"
)

type WeatherLocationService interface {
	GetLocation(ctx context.Context, id uint) (*dto.LocationResponse, error)
	ListLocations(ctx context.Context) (*dto.ListLocationResponse, error)
	CreateLocation(ctx context.Context, req dto.CreateWeatherLocationRequest) (*dto.LocationResponse, error)
	GetRainfall(ctx context.Context, id uint) ([]dto.WeatherRainfallResponse, error)
	GetTemperature(ctx context.Context, id uint) ([]dto.WeatherRecordResponse, error)
}

type weatherLocationService struct {
	locationRepo repository.WeatherLocationRepository
	tempRepo     repository.WeatherRecordRepository
	rainfallRepo repository.WeatherRainfallRepository
}

func NewWeatherLocationService(locationRepo repository.WeatherLocationRepository, tempRepo repository.WeatherRecordRepository, rainfallRepo repository.WeatherRainfallRepository) WeatherLocationService {
	return &weatherLocationService{
		locationRepo: locationRepo,
		tempRepo:     tempRepo,
		rainfallRepo: rainfallRepo,
	}
}

// GetLocation - Get a single location by ID
func (s *weatherLocationService) GetLocation(ctx context.Context, id uint) (*dto.LocationResponse, error) {
	location, err := s.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	temperature, err := s.tempRepo.List(ctx, location.ID)
	if err != nil {
		return nil, err
	}
	rainfall, err := s.rainfallRepo.List(ctx, location.ID)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(location, temperature, rainfall), nil
}

// ListLocations - Get all locations
func (s *weatherLocationService) ListLocations(ctx context.Context) (*dto.ListLocationResponse, error) {
	locations, err := s.locationRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}

	// Convert models to DTOs
	responses := make([]dto.LocationResponse, len(locations))
	for i, location := range locations {
		temperature, err := s.tempRepo.List(ctx, location.ID)
		if err != nil {
			return nil, err
		}
		rainfall, err := s.rainfallRepo.List(ctx, location.ID)
		if err != nil {
			return nil, err
		}
		responses[i] = *s.modelToResponse(&location, temperature, rainfall)
	}

	return &dto.ListLocationResponse{
		Locations: responses,
		Count:     len(responses),
	}, nil
}

// CreateLocation - Create a new location
func (s *weatherLocationService) CreateLocation(ctx context.Context, req dto.CreateWeatherLocationRequest) (*dto.LocationResponse, error) {
	// Business Rule: You could validate coordinates here
	if req.Latitude < -90 || req.Latitude > 90 {
		return nil, errors.New("latitude must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return nil, errors.New("longitude must be between -180 and 180")
	}

	// Create model
	location := &models.WeatherLocation{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// Save to database
	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return s.modelToResponse(location, nil, nil), nil
}

func modelTemp(t []models.WeatherRecord) []dto.WeatherRecordResponse {
	list := make([]dto.WeatherRecordResponse, len(t))
	for i := range t {
		list[i] = dto.WeatherRecordResponse{
			ID:            t[i].ID,
			Time:          t[i].Time,
			Temperature:   t[i].Temperature,
			Windspeed:     t[i].Windspeed,
			Winddirection: t[i].Winddirection,
		}
	}
	return list
}

func modelRailfall(t []models.Rainfall) []dto.WeatherRainfallResponse {
	list := make([]dto.WeatherRainfallResponse, len(t))
	for i := range t {
		list[i] = dto.WeatherRainfallResponse{
			ID:                 t[i].ID,
			Time:               t[i].Time,
			Precipitation:      t[i].Precipitation,
			Evapotranspiration: t[i].Evapotranspiration,
		}
	}
	return list
}

// GetLocation - Get a single location by ID
func (s *weatherLocationService) GetRainfall(ctx context.Context, id uint) ([]dto.WeatherRainfallResponse, error) {
	location, err := s.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	rainfall, err := s.rainfallRepo.List(ctx, location.ID)
	if err != nil {
		return nil, err
	}

	return modelRailfall(rainfall), nil
}

// GetLocation - Get a single location by ID
func (s *weatherLocationService) GetTemperature(ctx context.Context, id uint) ([]dto.WeatherRecordResponse, error) {
	location, err := s.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	temperature, err := s.tempRepo.List(ctx, location.ID)
	if err != nil {
		return nil, err
	}

	return modelTemp(temperature), nil
}

// Helper: Convert model to DTO
func (s *weatherLocationService) modelToResponse(location *models.WeatherLocation, temp []models.WeatherRecord, rainfall []models.Rainfall) *dto.LocationResponse {

	return &dto.LocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Temperature: modelTemp(temp),
		Rainfall:    modelRailfall(rainfall),
	}
}
