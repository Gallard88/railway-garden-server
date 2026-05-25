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
	GetLocation(ctx context.Context, id uint) (*dto.WeatherLocationResponse, error)
	ListLocations(ctx context.Context) (*dto.ListWeatherLocationsResponse, error)
	CreateLocation(ctx context.Context, req dto.CreateWeatherLocationRequest) (*dto.WeatherLocationResponse, error)
}

type weatherLocationService struct {
	locationRepo repository.WeatherLocationRepository
}

func NewWeatherLocationService(locationRepo repository.WeatherLocationRepository) WeatherLocationService {
	return &weatherLocationService{
		locationRepo: locationRepo,
	}
}

// GetLocation - Get a single location by ID
func (s *weatherLocationService) GetLocation(ctx context.Context, id uint) (*dto.WeatherLocationResponse, error) {
	location, err := s.locationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(location), nil
}

// ListLocations - Get all locations
func (s *weatherLocationService) ListLocations(ctx context.Context) (*dto.ListWeatherLocationsResponse, error) {
	locations, err := s.locationRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}

	// Convert models to DTOs
	responses := make([]dto.WeatherLocationResponse, len(locations))
	for i, location := range locations {
		responses[i] = *s.modelToResponse(&location)
	}

	return &dto.ListWeatherLocationsResponse{
		Locations: responses,
		Count:     len(responses),
	}, nil
}

// CreateLocation - Create a new location
func (s *weatherLocationService) CreateLocation(ctx context.Context, req dto.CreateWeatherLocationRequest) (*dto.WeatherLocationResponse, error) {
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

	return s.modelToResponse(location), nil
}

// Helper: Convert model to DTO
func (s *weatherLocationService) modelToResponse(location *models.WeatherLocation) *dto.WeatherLocationResponse {
	return &dto.WeatherLocationResponse{
		ID:        location.ID,
		Name:      location.Name,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		CreatedAt: location.CreatedAt,
	}
}
