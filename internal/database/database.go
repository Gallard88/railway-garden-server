package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"goapi.railway.app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() (*gorm.DB, error) {
	// Get database URL from environment
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("DATABASE_URL"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)
	//databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)
	if os.Getenv("GIN_MODE") == "release" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Successfully connected to database and ran migrations")
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	// Auto-migrate your models
	return db.AutoMigrate(
		&models.Plant{},
		&models.WeatherLocation{},
		&models.Weather{},
		&models.Rainfall{},
		// Add more models here as you create them
	)
}
