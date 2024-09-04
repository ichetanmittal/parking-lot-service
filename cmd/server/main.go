package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/chetan/parking-lot-service/internal/db"
	"github.com/chetan/parking-lot-service/internal/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database := db.ConnectDatabase()

	e := echo.New()

	// Initialize handlers
	h := handlers.NewHandler(database)

	// Define routes
	e.POST("/parking-lots", h.CreateParkingLot)
	e.GET("/parking-lots/:id", h.GetParkingLot)
	e.GET("/parking-lots/:id/available-spots", h.GetAvailableSpots)
	e.POST("/parking-entries", h.CreateParkingEntry)
	e.PUT("/parking-entries/:id/exit", h.ExitParking)
	e.POST("/tariffs", h.CreateTariff)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}