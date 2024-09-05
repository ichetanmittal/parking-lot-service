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

type Handler struct {
	DB             *gorm.DB
	ParkingService *services.ParkingService
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB:             db,
		ParkingService: services.NewParkingService(db),
	}
}

func (h *Handler) CreateParkingLot(c echo.Context) error {
	parkingLot := new(models.ParkingLot)
	if err := c.Bind(parkingLot); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := parkingLot.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.DB.Create(parkingLot).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create parking lot"})
	}

	return c.JSON(http.StatusCreated, parkingLot)
}

func (h *Handler) GetParkingLot(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	parkingLot := new(models.ParkingLot)

	if err := h.DB.First(parkingLot, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Parking lot not found"})
	}

	return c.JSON(http.StatusOK, parkingLot)
}

func (h *Handler) GetAvailableSpots(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid parking lot ID"})
	}

	availableSpots, err := h.ParkingService.GetAvailableSpots(uint(id))
	if err != nil {
		log.Printf("Error getting available spots: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get available spots"})
	}

	return c.JSON(http.StatusOK, availableSpots)
}

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

	response := struct {
		*models.ParkingEntry
		*models.Receipt
	}{
		ParkingEntry: entry,
		Receipt:      receipt,
	}

	return c.JSON(http.StatusOK, response)
}

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