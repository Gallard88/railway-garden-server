package service

import (
	"context"
	"fmt"
	"time"

	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/models"
	"goapi.railway.app/internal/repository"
)

type PlantZoneService interface {
	GetZone(ctx context.Context, id uint) (*dto.PlantZoneResponse, error)
	ListZones(ctx context.Context) (*dto.ListPlantZonesResponse, error)
	CreateZone(ctx context.Context, req dto.CreatePlantZoneRequest) (*dto.PlantZoneResponse, error)
}

type plantZoneService struct {
	zoneRepo repository.PlantZoneRepository
}

func NewPlantZoneService(zoneRepo repository.PlantZoneRepository) PlantZoneService {
	return &plantZoneService{
		zoneRepo: zoneRepo,
	}
}

// GetZone - Get a single zone by ID
func (s *plantZoneService) GetZone(ctx context.Context, id uint) (*dto.PlantZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(zone), nil
}

// ListZones - Get all zones
func (s *plantZoneService) ListZones(ctx context.Context) (*dto.ListPlantZonesResponse, error) {
	zones, err := s.zoneRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch zones: %w", err)
	}

	// Convert models to DTOs
	responses := make([]dto.PlantZoneResponse, len(zones))
	for i, zone := range zones {
		responses[i] = *s.modelToResponse(&zone)
	}

	return &dto.ListPlantZonesResponse{
		Zones: responses,
		Count: len(responses),
	}, nil
}

// CreateZone - Create a new zone
func (s *plantZoneService) CreateZone(ctx context.Context, req dto.CreatePlantZoneRequest) (*dto.PlantZoneResponse, error) {
	zone := &models.PlantZone{
		Name:          req.Name,
		LocationID:    req.LocationID,
		Outdoor:       req.Outdoor,
		RainThreshold: req.RainThreshold,
		CreatedAt:     time.Now(),
	}

	// Save to database
	if err := s.zoneRepo.Create(ctx, zone); err != nil {
		return nil, fmt.Errorf("failed to create zone: %w", err)
	}

	return s.modelToResponse(zone), nil
}

// Helper: Convert model to DTO
func (s *plantZoneService) modelToResponse(zone *models.PlantZone) *dto.PlantZoneResponse {
	return &dto.PlantZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		LocationID:    zone.LocationID,
		Outdoor:       zone.Outdoor,
		RainThreshold: zone.RainThreshold,
		CreatedAt:     zone.CreatedAt,
	}
}
