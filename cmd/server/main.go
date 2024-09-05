package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/chetan/parking-lot-service/internal/db"
	"github.com/chetan/parking-lot-service/internal/handlers"
)

// main is the entry point of the application.
// It sets up the environment, initializes the database connection,
// creates an Echo instance, sets up routes, and starts the server.
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	database := db.ConnectDatabase()

	// Create a new Echo instance
	e := echo.New()

	// Initialize handlers with the database connection
	h := handlers.NewHandler(database)

	// Define routes
	e.POST("/parking-lots", h.CreateParkingLot)
	e.GET("/parking-lots/:id", h.GetParkingLot)
	e.GET("/parking-lots/:id/available-spots", h.GetAvailableSpots)
	e.POST("/parking-entries", h.CreateParkingEntry)
	e.PUT("/parking-entries/:id/exit", h.ExitParking)
	e.POST("/tariffs", h.CreateTariff)

	// Get the port from environment variables, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	e.Logger.Fatal(e.Start(":" + port))
}