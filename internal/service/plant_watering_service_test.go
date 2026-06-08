package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"goapi.railway.app/internal/models"
)

// Mock repositories for testing
type MockPlantRepository struct {
	mock.Mock
}

func (m *MockPlantRepository) FindByZone(ctx context.Context, zoneID uint) ([]models.Plant, error) {
	args := m.Called(ctx, zoneID)
	if args.Get(0) == nil {
		return []models.Plant{}, args.Error(1)
	}
	return args.Get(0).([]models.Plant), args.Error(1)
}

func (m *MockPlantRepository) Water(ctx context.Context, id uint) (*models.Plant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) FindByID(ctx context.Context, id uint) (*models.Plant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plant), args.Error(1)
}

func (m *MockPlantRepository) List(ctx context.Context) ([]models.Plant, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []models.Plant{}, args.Error(1)
	}
	return args.Get(0).([]models.Plant), args.Error(1)
}

func (m *MockPlantRepository) Update(ctx context.Context, updatedPlanr models.Plant) error {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *MockPlantRepository) Create(ctx context.Context, plant *models.Plant) error {
	args := m.Called(ctx, plant)
	return args.Error(0)
}
func (m *MockPlantRepository) DeletePlant(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRainfallRepository struct {
	mock.Mock
}

func (m *MockRainfallRepository) List(ctx context.Context, id uint) ([]models.Rainfall, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]models.Rainfall), args.Error(1)
}

func (m *MockRainfallRepository) RainfallForLocationAndDate(ctx context.Context, locationId uint, date time.Time) (models.Rainfall, error) {
	args := m.Called(ctx, locationId, date)
	if args.Get(0) == nil {
		return models.Rainfall{}, args.Error(1)
	}
	return args.Get(0).(models.Rainfall), args.Error(1)
}

func (m *MockRainfallRepository) FindByZoneAndTime(ctx context.Context, zoneID uint, from, to time.Time) ([]models.Rainfall, error) {
	args := m.Called(ctx, zoneID, from, to)
	if args.Get(0) == nil {
		return []models.Rainfall{}, args.Error(1)
	}
	return args.Get(0).([]models.Rainfall), args.Error(1)
}

func (m *MockRainfallRepository) GetTotalByZoneAndTime(ctx context.Context, zoneID uint, from, to time.Time) (float64, error) {
	args := m.Called(ctx, zoneID, from, to)
	return args.Get(0).(float64), args.Error(1)
}

type MockPlantZoneRepository struct {
	mock.Mock
}

func (m *MockPlantZoneRepository) Create(ctx context.Context, event *models.PlantZone) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockPlantZoneRepository) FindByID(ctx context.Context, id uint) (*models.PlantZone, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlantZone), args.Error(1)
}
func (m *MockPlantZoneRepository) List(ctx context.Context) ([]models.PlantZone, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []models.PlantZone{}, args.Error(1)
	}
	return args.Get(0).([]models.PlantZone), args.Error(1)
}
func (m *MockPlantZoneRepository) ListExposedToRainfall(ctx context.Context) ([]models.PlantZone, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []models.PlantZone{}, args.Error(1)
	}
	return args.Get(0).([]models.PlantZone), args.Error(1)
}

// ============================================
// TESTS
// ============================================
/*

func TestProcessRainfall_RainfallBelowThreshold(t *testing.T) {
	// Setup mocks
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockPlantZones := new(MockPlantZoneRepository)
	// Expect rainfall to be saved
	mockRainfallRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Expect total rainfall to be 10mm (below 20mm threshold)
	mockRainfallRepo.On("GetTotalByZoneAndTime", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(10.0, nil)

	// No plants should be watered
	mockPlantRepo.AssertNotCalled(t, "FindByZone")

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockPlantZones)

	event := &models.Rainfall{
		ZoneID:     1,
		Precipitation:     10.0,
		RecordedAt: time.Now(),
	}

	result, err := service.ProcessRainfall(context.Background(), event)

	assert.NoError(t, err)
	assert.False(t, result.ThresholdMet)
	assert.Equal(t, 10.0, result.TotalRainfall)
	assert.Equal(t, 0, result.PlantsWatered)
}
*/
/*

func TestProcessRainfall_RainfallAboveThreshold(t *testing.T) {
	// Setup mocks
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockPlantZones := new(MockPlantZoneRepository)

	// Expect rainfall to be saved
	mockRainfallRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Expect total rainfall to be 25mm (above 20mm threshold)
	mockRainfallRepo.On("GetTotalByZoneAndTime", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(25.0, nil)

	// Should find plants in zone
	plants := []models.Plant{
		{ID: 1, Name: "Rose", ZoneID: 1},
		{ID: 2, Name: "Daisy", ZoneID: 1},
		{ID: 3, Name: "Tulip", ZoneID: 1},
	}
	mockPlantRepo.On("FindByZone", mock.Anything, uint(1)).Return(plants, nil)

	// Should update last_watered for each plant
	mockPlantRepo.On("UpdateLastWatered", mock.Anything, uint(1)).Return(nil)
	mockPlantRepo.On("UpdateLastWatered", mock.Anything, uint(2)).Return(nil)
	mockPlantRepo.On("UpdateLastWatered", mock.Anything, uint(3)).Return(nil)

	// Should log each watering
	mockWateringLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockWateringLogRepo)

	event := &models.Rainfall{
		ZoneID:     1,
		Precipitation:     25.0,
		RecordedAt: time.Now(),
	}

	result, err := service.ProcessRainfall(context.Background(), event)

	assert.NoError(t, err)
	assert.True(t, result.ThresholdMet)
	assert.Equal(t, 25.0, result.TotalRainfall)
	assert.Equal(t, 3, result.PlantsWatered)
	assert.Equal(t, []uint{1, 2, 3}, result.WateredPlants)
}
*/
/*

func TestProcessRainfall_NegativeRainfall(t *testing.T) {
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockWateringLogRepo := new(MockWateringLogRepository)

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockWateringLogRepo)

	event := &models.Rainfall{
		ZoneID:     1,
		Precipitation:     -5.0, // Invalid!
		RecordedAt: time.Now(),
	}

	result, err := service.ProcessRainfall(context.Background(), event)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot be negative")
}
*/
/*
func TestProcessRainfall_DatabaseError(t *testing.T) {
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockWateringLogRepo := new(MockWateringLogRepository)

	// Simulate database error
	mockRainfallRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("database connection failed"))

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockWateringLogRepo)

	event := &models.Rainfall{
		ZoneID:     1,
		Precipitation:     25.0,
		RecordedAt: time.Now(),
	}

	result, err := service.ProcessRainfall(context.Background(), event)

	assert.Error(t, err)
	assert.Nil(t, result)
}
*/
/*
func TestProcessRainfall_EdgeCase_ExactlyAtThreshold(t *testing.T) {
	// Test boundary condition: exactly 20mm
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockWateringLogRepo := new(MockWateringLogRepository)

	mockRainfallRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockRainfallRepo.On("GetTotalByZoneAndTime", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(20.0, nil)

	plants := []models.Plant{{ID: 1, Name: "Plant", ZoneID: 1}}
	mockPlantRepo.On("FindByZone", mock.Anything, uint(1)).Return(plants, nil)
	mockPlantRepo.On("UpdateLastWatered", mock.Anything, mock.Anything).Return(nil)
	mockWateringLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockWateringLogRepo)

	event := &models.Rainfall{
		ZoneID:     1,
		Precipitation:     20.0,
		RecordedAt: time.Now(),
	}

	result, err := service.ProcessRainfall(context.Background(), event)

	assert.NoError(t, err)
	assert.True(t, result.ThresholdMet) // Should be watered at exactly 20mm
	assert.Equal(t, 1, result.PlantsWatered)
}
*/

func TestProcessRainfall_NoPlants(t *testing.T) {
	// Test when zone has no plants
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockPlantZoneRepository := new(MockPlantZoneRepository)

	mockRainfallRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockRainfallRepo.On("GetTotalByZoneAndTime", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(25.0, nil)
	mockPlantRepo.On("FindByZone", mock.Anything, uint(1)).Return([]models.Plant{}, nil)

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockPlantZoneRepository)

	result, err := service.ProcessRainfallEvent(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 0, result.PlantsWatered)
	assert.Equal(t, 0, len(result.WateredPlants))
}

/*
func TestRecordRainfall_HigherLevelMethod(t *testing.T) {
	// Test the convenience method
	mockPlantRepo := new(MockPlantRepository)
	mockRainfallRepo := new(MockRainfallRepository)
	mockWateringLogRepo := new(MockWateringLogRepository)

	mockRainfallRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *models.Rainfall) bool {
		return e.ZoneID == 2 && e.Precipitation == 30.0
	})).Return(nil)

	mockRainfallRepo.On("GetTotalByZoneAndTime", mock.Anything, uint(2), mock.Anything, mock.Anything).Return(30.0, nil)

	plants := []models.Plant{{ID: 10, Name: "Lily", ZoneID: 2}}
	mockPlantRepo.On("FindByZone", mock.Anything, uint(2)).Return(plants, nil)
	mockPlantRepo.On("UpdateLastWatered", mock.Anything, uint(10)).Return(nil)
	mockWateringLogRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	service := NewPlantWateringService(mockPlantRepo, mockRainfallRepo, mockWateringLogRepo)

	result, err := service.RecordRainfall(context.Background(), 2, 30.0, time.Now())

	assert.NoError(t, err)
	assert.True(t, result.ThresholdMet)
	assert.Equal(t, 1, result.PlantsWatered)
}
*/
