package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/chetan/parking-lot-service/internal/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_ = db.ConnectDatabase()

	e := echo.New()

	// TODO: Add routes and handlers

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}