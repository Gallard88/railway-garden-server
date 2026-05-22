package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"goapi.railway.app/internal/database"
)

const version = "0.0.1"

type config struct {
	port   int
	env    string
	db_url string
}

type application struct {
	config *config
	logger *slog.Logger
	db     *gorm.DB
	cron   *cron.Cron
}

func main() {
	cfg := Load()
	// Connect to database
	// Connect to database
	db, err := database.New()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// create the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize cron scheduler
	cronScheduler := cron.New(cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
	))
	// create the application
	app := &application{
		config: cfg,
		logger: logger,
		cron:   cronScheduler,
		db:     db,
	}

	// Setup cron
	app.setupCronJobs()
	app.cron.Start()
	defer app.cron.Stop()
	log.Println("Cron scheduler started")

	// Set Gin mode
	gin.SetMode(cfg.env)

	// Create router and start server
	router := app.routes()

	addr := fmt.Sprintf(":%d", cfg.port)
	log.Printf("Server starting on %s (version %s)", addr, version)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func Load() *config {
	// Load .env file if it exists (optional in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "postgres"),
	)
	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 8080
	}
	return &config{
		db_url: dbURL,
		port:   intPort,
		env:    getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Run every 15 minutes:
//app.cron.AddFunc("*/15 * * * *", services.FetchWeatherForAllLocations)
//services.FetchWeatherForAllLocations()
//app.cron.AddFunc("* 20 * * *", services.FetchRainfallForAllLocations)
//services.FetchRainfallForAllLocations()

//app.cron.Start()
// Setup cron jobs
//	app.setupCronJobs()

// Start cron scheduler
//	app.cron.Start()
//	defer app.cron.Stop()
