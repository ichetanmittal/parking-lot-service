package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chetan/parking-lot-service/internal/errors"
	"github.com/chetan/parking-lot-service/internal/models"
	"github.com/chetan/parking-lot-service/internal/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler struct holds dependencies for the HTTP handlers
type Handler struct {
	DB             *gorm.DB
	ParkingService *services.ParkingService
}

// NewHandler creates a new Handler instance with the given database connection
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB:             db,
		ParkingService: services.NewParkingService(db),
	}
}

// CreateParkingLot handles the creation of a new parking lot
func (h *Handler) CreateParkingLot(c echo.Context) error {
	parkingLot := new(models.ParkingLot)
	// Uses Echo's Bind method to parse the request body into the parkingLot struct.
	// If the binding fails, it returns a 400 Bad Request response with an error message.
	if err := c.Bind(parkingLot); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Validates the parkingLot struct using the Validate method (models package).
	// If the validation fails, it returns a 400 Bad Request response with an error message.
	if err := parkingLot.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Creates a new parking lot in the database using GORM's Create method.
	// If the creation fails, it returns a 500 Internal Server Error response with an error message.

	if err := h.DB.Create(parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create parking lot"})
	}

	// Returns a 201 Created response with the newly created parking lot data.
	return c.JSON(http.StatusCreated, parkingLot)
}

// GetParkingLot retrieves details of a specific parking lot
func (h *Handler) GetParkingLot(c echo.Context) error {
	// Parses the parking lot ID from the URL parameter and converts it to an integer.
	id, _ := strconv.Atoi(c.Param("id"))
	parkingLot := new(models.ParkingLot)

	// If the retrieval fails, it returns a 404 Not Found response with an error message.
	if err := h.DB.First(parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Parking lot not found"})
	}

	// Returns a 200 OK response with the parking lot data.
	return c.JSON(http.StatusOK, parkingLot)
}

// GetAvailableSpots retrieves the number of available spots for each vehicle type in a parking lot
func (h *Handler) GetAvailableSpots(c echo.Context) error {
	// Parses the parking lot ID from the URL parameter and converts it to an integer.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid parking lot ID"})
	}

	// Retrieves the available spots for the specified parking lot using the ParkingService.
	availableSpots, err := h.ParkingService.GetAvailableSpots(uint(id))
	if err != nil {
		log.Printf("Error getting available spots: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get available spots"})
	}

	return c.JSON(http.StatusOK, availableSpots)
}

// CreateParkingEntry handles the creation of a new parking entry
func (h *Handler) CreateParkingEntry(c echo.Context) error {

	entry := new(models.ParkingEntry)
	if err := c.Bind(entry); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidInput.Error()})
	}

	if err := h.ParkingService.CreateParkingEntry(entry); err != nil {
		switch err {
		case errors.ErrParkingLotNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case errors.ErrNoAvailableSpots:
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create parking entry"})
		}
	}

	return c.JSON(http.StatusCreated, entry)
}

// ExitParking handles the process of a vehicle exiting the parking lot
func (h *Handler) ExitParking(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidInput.Error()})
	}

	entry, receipt, err := h.ParkingService.ExitParking(uint(id))
	if err != nil {
		switch err {
		case errors.ErrParkingEntryNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case errors.ErrVehicleAlreadyExited:
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		case errors.ErrTariffNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tariff not found for this parking lot and vehicle type"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process parking exit"})
		}
	}

	// Constructs a response struct containing both the parking entry and receipt data.
	// The purpose of this structure is to combine the ParkingEntry and Receipt into a single object for the JSON response. When this struct is serialized to JSON, it will create a nested structure where the fields of ParkingEntry and Receipt are at the top level of the JSON object.
	response := struct {
		*models.ParkingEntry
		*models.Receipt
	}{
		ParkingEntry: entry,
		Receipt:      receipt,
	}

	return c.JSON(http.StatusOK, response)
}

// CreateTariff handles the creation of a new tariff
func (h *Handler) CreateTariff(c echo.Context) error {
	tariff := new(models.Tariff)
	if err := c.Bind(tariff); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := tariff.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.ParkingService.CreateTariff(tariff)
	if err != nil {
		switch err {
		case errors.ErrParkingLotNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Parking lot not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tariff"})
		}
	}

	return c.JSON(http.StatusCreated, tariff)
}