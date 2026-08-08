package service

import (
	"context"
	"fmt"

	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/models"
	"goapi.railway.app/internal/repository"
)

type PlantHistoryService interface {
	ReadHistory(ctx context.Context, id uint) (*dto.PlantHistoryResponse, error)
	CreateHistoryEvent(ctx context.Context, event dto.WriteHistoryRequest) (bool, error)
}

type plantHistoryService struct {
	historyRepo repository.PlantHistoryRepository
}

func NewPlantHistoryService(historyRepo repository.PlantHistoryRepository) PlantHistoryService {
	return &plantHistoryService{
		historyRepo: historyRepo,
	}
}

//	ReadHistory(ctx context.Context, params dto.ReadHistoryRequest) (*dto.PlantHistoryResponse, error)
//	CreateHistoryEvent(ctx context.Context, event dto.WriteHistoryRequest) error

// ReadHistory - Get a single zone by ID
func (s *plantHistoryService) ReadHistory(ctx context.Context, id uint) (*dto.PlantHistoryResponse, error) {
	history, err := s.historyRepo.ReadHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.PlantHistoryResponse{Log: history}, nil
}

// CreateZone - Create a new zone
func (s *plantHistoryService) CreateHistoryEvent(ctx context.Context, req dto.WriteHistoryRequest) (bool, error) {
	entry := &models.PlantHistory{
		PlantId:     req.PlantId,
		Name:        req.Name,
		Description: req.Description,
		Agent:       req.Agent,
	}

	// Save to database
	if err := s.historyRepo.CreateHistoryEvent(ctx, entry); err != nil {
		return false, fmt.Errorf("failed to create hitory entry: %w", err)
	}

	return true, nil
}
