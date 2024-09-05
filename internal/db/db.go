package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/chetan/parking-lot-service/internal/models"
)

// ConnectDatabase establishes a connection to the PostgreSQL database
// and performs auto-migration for the defined models.
//
// Returns:
//   - *gorm.DB: A pointer to the initialized database connection.
func ConnectDatabase() *gorm.DB {
	// Get database connection details from environment variables
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "host.docker.internal" // Default to Docker host
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Perform auto-migration for the defined models
	err = db.AutoMigrate(&models.ParkingLot{}, &models.Tariff{}, &models.ParkingEntry{}, &models.Receipt{})
	if err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
	}

	return db
}
