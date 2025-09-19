package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sonic-labs/course-enrollment-service/internal/config"
	"sonic-labs/course-enrollment-service/internal/database"
	"sonic-labs/course-enrollment-service/internal/router"

	_ "sonic-labs/course-enrollment-service/docs" // Import generated docs
)

// @title Course Enrollment Service API
// @version 1.0
// @description A robust backend service for managing course catalog and student enrollments for Sonic University platform
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.sonic-labs.com/support
// @contact.email support@sonic-labs.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func setupLogging() {
	// Create logs directory
	logsDir := "/app/logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		return
	}

	// Create log file
	logFile := filepath.Join(logsDir, "app.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}

	// Set log output to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Logging setup completed - writing to file and stdout")
}

func main() {
	// Setup logging to file
	setupLogging()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed database with demo data
	if err := database.Seed(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Setup router
	r := router.Setup(db, cfg)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
