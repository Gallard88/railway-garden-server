package service

import (
	"context"
	"fmt"

	"goapi.railway.app/internal/dto"
	"goapi.railway.app/internal/models"
	"goapi.railway.app/internal/repository"
)

type WeatherRecordService interface {
	ListRecords(ctx context.Context, id uint) (*dto.ListWeatherRecordsResponse, error)
}

type weatherRecordService struct {
	recordRepo repository.WeatherRecordRepository
}

func NewWeatherRecordService(recordRepo repository.WeatherRecordRepository) WeatherRecordService {
	return &weatherRecordService{
		recordRepo: recordRepo,
	}
}

// ListRecords - Get all records
func (s *weatherRecordService) ListRecords(ctx context.Context, id uint) (*dto.ListWeatherRecordsResponse, error) {
	records, err := s.recordRepo.List(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}

	// Convert models to DTOs
	responses := make([]dto.WeatherRecordResponse, len(records))
	for i, record := range records {
		responses[i] = *s.modelToResponse(&record)
	}

	return &dto.ListWeatherRecordsResponse{
		Records: responses,
		Count:   len(responses),
	}, nil
}

// Helper: Convert model to DTO
func (s *weatherRecordService) modelToResponse(record *models.WeatherRecord) *dto.WeatherRecordResponse {
	return &dto.WeatherRecordResponse{
		ID: record.ID,
		//Name:      record.Name,
		LocationID:    record.LocationID,
		Time:          record.Time,
		Temperature:   record.Temperature,
		Windspeed:     record.Windspeed,
		Winddirection: record.Winddirection,
		CreatedAt:     record.CreatedAt,
	}
}
