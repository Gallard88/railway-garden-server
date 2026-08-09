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

// ================================================
type PlantService interface {
	GetPlant(ctx context.Context, id uint) (*dto.PlantResponse, error)
	UpdatePlant(ctx context.Context, id uint, plant models.Plant) (*dto.PlantResponse, error)
	ListPlants(ctx context.Context) (*dto.ListPlantsResponse, error)
	CreatePlant(ctx context.Context, req dto.CreatePlantRequest) (*dto.PlantResponse, error)
	DeletePlant(ctx context.Context, id uint) error
	Water(ctx context.Context, id uint) (*dto.PlantResponse, error)

	ReadHistory(ctx context.Context, id uint) (*dto.PlantHistoryResponse, error)
	CreateHistoryEvent(ctx context.Context, event dto.WriteHistoryRequest) (bool, error)
}

type plantService struct {
	plantRepo repository.PlantRepository
}

func NewPlantService(plantRepo repository.PlantRepository) PlantService {
	return &plantService{
		plantRepo: plantRepo,
	}
}

// GetPlant - Get a single plant by ID
func (s *plantService) GetPlant(ctx context.Context, id uint) (*dto.PlantResponse, error) {
	plant, err := s.plantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(plant), nil
}

func (s *plantService) UpdatePlant(ctx context.Context, id uint, updatedPlant models.Plant) (*dto.PlantResponse, error) {

	// Get Plant
	plant, err := s.plantRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Business rule validation
	if updatedPlant.ContainerType != models.ContainerType_Pot && updatedPlant.ContainerType != models.ContainerType_Ground {
		return nil, fmt.Errorf("containter type must be (%s,%s)", models.ContainerType_Pot, models.ContainerType_Ground)
	}
	if updatedPlant.SunExposure != models.Sun_FullSun && updatedPlant.SunExposure != models.Sun_PartSun && updatedPlant.SunExposure != models.Sun_FullShade {
		return nil, fmt.Errorf("Sun Exposure type must be (%s,%s,%s)", models.Sun_FullSun, models.Sun_PartSun, models.Sun_FullShade)
	}
	if updatedPlant.SoilType != models.Soil_Sandy && updatedPlant.SoilType != models.Soil_Clay && updatedPlant.SoilType != models.Soil_Loam {
		return nil, fmt.Errorf("SoilType type must be (%s,%s,%s)", models.Soil_Sandy, models.Soil_Clay, models.Soil_Loam)
	}
	if updatedPlant.ET0 <= 0 || updatedPlant.ET0 > 1.0 {
		return nil, fmt.Errorf("ETO0 must be between 0 & 1")
	}
	if updatedPlant.DeficitThreshold <= 0 || updatedPlant.DeficitThreshold > 100.0 {
		return nil, fmt.Errorf("DeficitThreshold must be between 0 & 100 mm")
	}
	if updatedPlant.LookbackDays < 1 || updatedPlant.LookbackDays > 7 {
		return nil, fmt.Errorf("LookbackDays must be between 1 & 7 days")
	}
	if updatedPlant.RainfallEffectiveness < 0 || updatedPlant.RainfallEffectiveness > 1.0 {
		return nil, fmt.Errorf("RainfallEffectiveness must be between 0 & 1")
	}

	// Update the appropriate fields.
	updatedPlant.CreatedAt = plant.CreatedAt
	updatedPlant.ID = plant.ID // to be sure.
	updatedPlant.UpdatedAt = time.Now()
	updatedPlant.LastWatered = plant.LastWatered
	updatedPlant.NextWater = plant.NextWater

	// Save to the DB.
	err = s.plantRepo.Update(ctx, updatedPlant)
	if err != nil {
		return nil, err
	}

	// Return
	return s.modelToResponse(&updatedPlant), nil
}

// Water - Mark a plant as watered
func (s *plantService) Water(ctx context.Context, id uint) (*dto.PlantResponse, error) {
	plant, err := s.plantRepo.Water(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.modelToResponse(plant), nil
}

// ListPlants - Get all plants
func (s *plantService) ListPlants(ctx context.Context) (*dto.ListPlantsResponse, error) {
	plants, err := s.plantRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plants: %w", err)
	}

	// Convert models to DTOs
	responses := make([]dto.PlantResponse, len(plants))
	for i, plant := range plants {
		responses[i] = *s.modelToResponse(&plant)
	}

	return &dto.ListPlantsResponse{
		Plants: responses,
		Count:  len(responses),
	}, nil
}

// CreatePlant - Create a new plant
func (s *plantService) CreatePlant(ctx context.Context, req dto.CreatePlantRequest) (*dto.PlantResponse, error) {
	plant := &models.Plant{
		Name:                  req.Name,
		Zone:                  req.Zone,
		WaterFreq:             req.WaterFreq,
		PlantedDate:           req.PlantedDate,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		LastWatered:           time.Now(),
		NextWater:             time.Now(),
		ContainerType:         models.ContainerType_Pot, // in-ground or pot
		SunExposure:           models.Sun_FullSun,       // exposed, part shade, full shade
		SoilType:              models.Soil_Loam,         // sandy, loam, clay
		ET0:                   1.0,                      // derived from sun exposure (1.0 / 0.6 / 0.3)
		DeficitThreshold:      5,                        // mm of deficit before watering triggered (lower for sandy, higher for clay)
		LookbackDays:          1,                        // 1 for pots, 3-5 for in-ground
		RainfallEffectiveness: 1.0,                      // how much rainfall actually reaches the zone (pots under eaves etc), default 1.0
	}
	// Save to database
	if err := s.plantRepo.Create(ctx, plant); err != nil {
		return nil, fmt.Errorf("failed to create plant: %w", err)
	}

	return s.modelToResponse(plant), nil
}

func (s *plantService) DeletePlant(ctx context.Context, id uint) error {
	return s.plantRepo.DeletePlant(ctx, id)
}

// ReadHistory - Get a single zone by ID
func (s *plantService) ReadHistory(ctx context.Context, id uint) (*dto.PlantHistoryResponse, error) {
	history, err := s.plantRepo.ReadHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PlantHistoryResponse{Log: history}, nil
}

// CreateHistoryEvent - Create a new zone
func (s *plantService) CreateHistoryEvent(ctx context.Context, req dto.WriteHistoryRequest) (bool, error) {
	entry := &models.PlantHistory{
		PlantId:     req.PlantId,
		Name:        req.Name,
		Description: req.Description,
		Agent:       req.Agent,
	}

	// Save to database
	if err := s.plantRepo.CreateHistoryEvent(ctx, entry); err != nil {
		return false, fmt.Errorf("failed to create hitory entry: %w", err)
	}

	return true, nil
}

// Helper: Convert model to DTO
func (s *plantService) modelToResponse(plant *models.Plant) *dto.PlantResponse {
	return &dto.PlantResponse{
		ID:                    plant.ID,
		Name:                  plant.Name,
		Zone:                  plant.Zone,
		WaterFreq:             plant.WaterFreq,
		PlantedDate:           plant.PlantedDate,
		CreatedAt:             plant.CreatedAt,
		UpdatedAt:             plant.UpdatedAt,
		LastWatered:           plant.LastWatered,
		NextWater:             plant.NextWater,
		ContainerType:         plant.ContainerType,
		SunExposure:           plant.SunExposure,
		SoilType:              plant.SoilType,
		ET0:                   plant.ET0,
		DeficitThreshold:      plant.DeficitThreshold,
		LookbackDays:          plant.LookbackDays,
		RainfallEffectiveness: plant.RainfallEffectiveness,
	}
}
