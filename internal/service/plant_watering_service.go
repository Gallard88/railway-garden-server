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
	CheckForOverheaterPlants(ctx context.Context)
	CheckWateringStatus(ctx context.Context)
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
	rainfall      repository.WeatherRainfallRepository
}

func NewPlantWateringService(
	plantRepo repository.PlantRepository,
	rainfallRepo repository.WeatherRainfallRepository,
	plantZones repository.PlantZoneRepository,
	weatherRecord repository.WeatherRecordRepository,
	rainfall repository.WeatherRainfallRepository,
) PlantWateringService {
	return &plantWateringService{
		plantRepo:     plantRepo,
		rainfallRepo:  rainfallRepo,
		plantZones:    plantZones,
		weatherRecord: weatherRecord,
		rainfall:      rainfall,
	}
}

func (s *plantWateringService) CheckForOverheaterPlants(ctx context.Context) {
	log.Println("Checking for over heated plants. ")
	/*
		Logic.
		For each weather location
			Get the last 24 hours of temp records.
			Grab the zone in that location, grab the plants.
			For each plant,
			drop records older than last watered.
			get max temp
			if temp > threshold.
				set over temp flag.
				create history entry.
	*/
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
				s.plantRepo.MarkNeedsWatering(ctx, plant.ID, "Overheated", fmt.Sprintf("Temp was %2.1f, threshold %2.1f ", maxTemp, plant.OverheatedTemp))
			}
		}
	}
}

func (s *plantWateringService) CheckWateringStatus(ctx context.Context) {
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
		// Get Plants in the zone.
		plants, err := s.plantRepo.FindByZoneID(ctx, zone.ID)
		if err != nil {
			log.Printf("Failed to get plants for zone %d: %s\n", zone.Name, err)
			return
		}
		for _, plant := range plants {
			pastTime := time.Now().Add(-time.Duration(plant.LookbackDays) * 24 * time.Hour)
			rainData, err := s.rainfallRepo.HistoricRainfall(ctx, zone.LocationID, pastTime)
			if err != nil {
				log.Printf("Failed to get temp for zone %d: %s\n", zone.Name, err)
				return
			}
			needWatering, msg := s.checkWateringStatus(plant, rainData)
			if needWatering {
				s.plantRepo.MarkNeedsWatering(ctx, plant.ID, "Evaporation", msg)
			}
		}
	}
}

func (s *plantWateringService) checkWateringStatus(plant models.Plant, list []models.Rainfall) (bool, string) {
	//deficit = sum over lookback_days of:
	//    (zone.transpiration * plant.et0_multiplier) - (zone.rainfall * plant.rainfall_effectiveness)
	// Set flag needs_watering
	transpiration := 0.0
	precipitation := 0.0

	for i := range list {
		transpiration += list[i].Evapotranspiration
		precipitation += list[i].Precipitation
	}
	adjusted_transpiration := transpiration * plant.ET0
	adjusted_precipitation := precipitation * plant.RainfallEffectiveness
	deficet := adjusted_transpiration - adjusted_precipitation

	log.Printf("%d:%s", plant.ID, plant.Name)
	log.Printf("Transpiration: %2.1f, %2.1f (adj)", transpiration, adjusted_transpiration)
	log.Printf("Precipitation: %2.1f, %2.1f (adj)", precipitation, adjusted_precipitation)
	log.Printf("Deficet: %2.1f, Threshold %2.1f", deficet, plant.DeficitThreshold)

	return deficet > plant.DeficitThreshold, fmt.Sprintf("Total transpiration: %2.1f, precipitation: %2.1f, threshold: %2.1f", transpiration, precipitation, plant.DeficitThreshold)
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
