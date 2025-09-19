package database

import (
	"fmt"
	"log"
	"sonic-labs/course-enrollment-service/internal/config"
	"sonic-labs/course-enrollment-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize creates a new database connection
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connection established successfully")
	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.Course{},
		&models.Enrollment{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func Seed(db *gorm.DB) error {
	log.Println("Seeding database with demo data...")

	var count int64
	db.Model(&models.Course{}).Count(&count)
	if count > 0 {
		log.Println("Demo courses already exist, skipping seeding")
		return nil
	}
	demoCourses := []models.Course{
		{
			Title:       "Introduction to Go Programming",
			Description: "Learn the fundamentals of Go programming language, including syntax, data types, functions, and basic concurrency patterns.",
			Difficulty:  "Beginner",
		},
		{
			Title:       "Advanced Web Development with React",
			Description: "Master advanced React concepts including hooks, context, state management, and building scalable web applications.",
			Difficulty:  "Advanced",
		},
		{
			Title:       "Database Design and SQL",
			Description: "Comprehensive course covering relational database design principles, SQL queries, indexing, and performance optimization.",
			Difficulty:  "Intermediate",
		},
		{
			Title:       "Machine Learning Fundamentals",
			Description: "Introduction to machine learning algorithms, data preprocessing, model training, and evaluation techniques using Python.",
			Difficulty:  "Intermediate",
		},
	}

	for _, course := range demoCourses {
		if err := db.Create(&course).Error; err != nil {
			return fmt.Errorf("failed to create demo course '%s': %w", course.Title, err)
		}
		log.Printf("Created demo course: %s", course.Title)
	}

	log.Println("Database seeding completed successfully")
	return nil
}
