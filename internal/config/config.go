package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Port      string         `mapstructure:"PORT"`
	Database  DatabaseConfig `mapstructure:"database"`
	JWTSecret string         `mapstructure:"JWT_SECRET"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// Load loads configuration from environment variables
func Load() *Config {
	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "course_enrollment")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("JWT_SECRET", "your-default-jwt-secret-change-this")
	viper.SetDefault("admin.username", "admin")
	viper.SetDefault("admin.password", "admin!dev")

	// Load from environment variables
	viper.AutomaticEnv()

	// Map environment variables to config structure
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("PORT", port)
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database.host", dbHost)
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		viper.Set("database.port", dbPort)
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		viper.Set("database.user", dbUser)
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database.password", dbPassword)
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database.dbname", dbName)
	}
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		viper.Set("database.sslmode", sslMode)
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		viper.Set("JWT_SECRET", jwtSecret)
	}
	if adminUsername := os.Getenv("ADMIN_USERNAME"); adminUsername != "" {
		viper.Set("admin.username", adminUsername)
	}
	if adminPassword := os.Getenv("ADMIN_PASSWORD"); adminPassword != "" {
		viper.Set("admin.password", adminPassword)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &config
}
