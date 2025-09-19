package database

import (
	"log"

	"sonic-labs/course-enrollment-service/internal/constants"
	"sonic-labs/course-enrollment-service/internal/models"
	"sonic-labs/course-enrollment-service/internal/service"

	"gorm.io/gorm"
)

// SeedAdminUser creates the default admin user if it doesn't exist
// Only runs after migration is complete
func SeedAdminUser(db *gorm.DB) error {
	// Simple check if admin user already exists
	var count int64
	err := db.Model(&models.User{}).Where("username = ?", "admin").Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping seed")
		return nil
	}

	// Hash the default password using bcrypt
	hashedPassword, err := service.HashPassword("admin!dev")
	if err != nil {
		return err
	}

	// Simple SQL insert for admin user
	adminUser := models.User{
		Username: "admin",
		Password: hashedPassword,
		Role:     constants.RoleAdmin,
	}

	err = db.Create(&adminUser).Error
	if err != nil {
		return err
	}

	log.Printf("Admin user created successfully with ID: %s", adminUser.ID.String())
	return nil
}
