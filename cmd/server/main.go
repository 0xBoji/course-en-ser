package main

import (
	"log"
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
func main() {
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
