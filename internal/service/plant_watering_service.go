package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"goapi.railway.app/internal/models"
	"goapi.railway.app/internal/repository"
)

const (
	// Business Rule: Minimum rainfall to mark plants as watered (in mm)
	RAINFALL_THRESHOLD = 20.0
)

type PlantWateringService interface {
	ProcessRainfallEvent(ctx context.Context) (*ProcessRainfallResult, error)
	CheckForOverheaterPlants(ctx context.Context)
}

// ProcessRainfallResult - What happened when we processed rainfall
type ProcessRainfallResult struct {
	PlantsWatered int    `json:"plants_watered"`
	WateredPlants []uint `json:"watered_plant_ids"`
}

type plantWateringService struct {
	plantRepo     repository.PlantRepository
	rainfallRepo  repository.WeatherRainfallRepository
	plantZones    repository.PlantZoneRepository
	weatherRecord repository.WeatherRecordRepository
}

func NewPlantWateringService(
	plantRepo repository.PlantRepository,
	rainfallRepo repository.WeatherRainfallRepository,
	plantZones repository.PlantZoneRepository,
	weatherRecord repository.WeatherRecordRepository,
) PlantWateringService {
	return &plantWateringService{
		plantRepo:     plantRepo,
		rainfallRepo:  rainfallRepo,
		plantZones:    plantZones,
		weatherRecord: weatherRecord,
	}
}

// ProcessRainfallEvent - Handle a single rainfall event
func (s *plantWateringService) ProcessRainfallEvent(ctx context.Context) (*ProcessRainfallResult, error) {

	fmt.Printf("=====================\nProcessing rainfall event at %s\n", time.Now().Format(time.RFC3339))
	defer fmt.Printf("=====================\n")
	// 1. Get a list of relevant zones.
	zones, err := s.plantZones.ListExposedToRainfall(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list plant zones: %w", err)
	}

	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	fmt.Printf("Time range: %s - %s\n", startOfDay.Format(time.RFC3339), endOfDay.Format(time.RFC3339))

	result := &ProcessRainfallResult{
		PlantsWatered: 0,
		WateredPlants: []uint{},
	}

	for _, zone := range zones {
		fmt.Printf("Processing rainfall event for zone: %+v\n", zone)
		rainfall, err := s.rainfallRepo.RainfallForLocationAndDate(ctx, zone.LocationID, startOfDay)
		if err != nil {
			return nil, fmt.Errorf("failed to get rainfall for location %d: %w", zone.LocationID, err)
		}
		fmt.Printf("Retrieved rainfall for zone %d: %+v\n", zone.LocationID, rainfall)
		// 2. For each zone get most recent rainfall and if above threshold, water plants in that zone.
		if rainfall.Precipitation >= zone.RainThreshold {
			fmt.Printf("Rainfall %.1fmm exceeds threshold (%.1fmm). Marking plants as watered.\n", rainfall.Precipitation, zone.RainThreshold)

			plantRepos, err := s.plantRepo.FindByZone(ctx, zone.ID)
			if err != nil {
				fmt.Printf("Failed to get plants for zone %d: %v\n", zone.ID, err)
				continue
			}
			plantIDs := make([]uint, len(plantRepos))
			for i, plant := range plantRepos {
				plantIDs[i] = plant.ID
				_, err := s.plantRepo.Water(ctx, plant.ID)
				if err != nil {
					fmt.Printf("Failed to water plant %d: %v\n", plant.ID, err)
					continue
				}
				result.PlantsWatered++
			}
			result.WateredPlants = append(result.WateredPlants, plantIDs...)
		}
	}

	return result, nil
}

/*
// waterPlantsInZone - Private helper to mark plants as watered

	func (s *plantWateringService) waterPlantsInZone(ctx context.Context, zoneID uint, rainfallAmount float64) ([]uint, error) {
		// 1. Get all plants in the zone
		plants, err := s.plantRepo.FindByZone(ctx, zoneID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plants: %w", err)
		}

		if len(plants) == 0 {
			return []uint{}, nil
		}

		wateredPlantIDs := []uint{}

		// 2. For each plant, update last_watered and log it
		for _, plant := range plants {
			// Update plant's last_watered timestamp
			if err := s.plantRepo.UpdateLastWatered(ctx, plant.ID); err != nil {
				log.Printf("Warning: Failed to update plant %d: %v", plant.ID, err)
				continue
			}

			// Log this watering event
			log.Entry := &models.PlantWateringLog{
				PlantID: plant.ID,
				Reason:  "rainfall",
				Amount:  rainfallAmount,
			}

			if err := s.wateringLogRepo.Create(ctx, log.Entry); err != nil {
				log.Printf("Warning: Failed to log watering for plant %d: %v", plant.ID, err)
				continue
			}

			wateredPlantIDs = append(wateredPlantIDs, plant.ID)
		}

		return wateredPlantIDs, nil
	}
*/
func (s *plantWateringService) CheckForOverheaterPlants(ctx context.Context) {
	log.Println("Checking for over heated plants. ")
	/*
		Logic.
		For each weather location
			Get the last 24 hours of temp records.
			* Sort oldest to newest.
			Grab the zone in that location.
			For each Zone, grab the plants.
			* This could probably be one super query.

			For each plant,
			drop records older than last watered.
			get max temp
			if temp > threshold.
				set over temp flag.
				create history entry.
	*/
	log.Printf("Check for over heated plants\n")
	zones, err := s.plantZones.List(ctx)
	if err != nil {
		log.Printf("Failed to get zones: %s\n", err)
		return
	}
	for _, zone := range zones {
		if zone.LocationID == 0 {
			log.Printf("Skip Zone: %s\n", zone.Name)
			continue
		}
		// Get Weather data
		temp, err := s.weatherRecord.ListByTime(ctx, zone.LocationID, time.Now().Add(-24*time.Hour))
		if err != nil {
			log.Printf("Failed to get temp for zone %d: %s\n", zone.Name, err)
			return
		}

		// Get Plants in the zone.
		plants, err := s.plantRepo.FindByZoneID(ctx, zone.ID)
		if err != nil {
			log.Printf("Failed to get plants for zone %d: %s\n", zone.Name, err)
			return
		}
		for _, plant := range plants {
			maxTemp := s.findMaxTempWithinDate(temp, plant.LastWatered)
			if maxTemp >= plant.OverheatedTemp {
				s.plantRepo.MarkOverheated(ctx, plant.ID, fmt.Sprintf("Temp was %2.1f, threshold %2.1f ", maxTemp, plant.OverheatedTemp))
			}
		}
	}
}

func (s *plantWateringService) findMaxTempWithinDate(list []models.WeatherRecord, date time.Time) float64 {
	maxTemp := list[0].Temperature
	for _, record := range list {
		if record.Time.Before(date) {
			continue
		}
		if record.Temperature > maxTemp {
			maxTemp = record.Temperature
		}
	}
	return maxTemp
}
