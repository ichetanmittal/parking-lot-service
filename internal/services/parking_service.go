package services

import (
	"math"
	"time"

	"github.com/chetan/parking-lot-service/internal/errors"
	"github.com/chetan/parking-lot-service/internal/models"
	"gorm.io/gorm"
)

// ParkingService handles business logic for parking-related operations.
// It encapsulates database operations and implements the core functionality
// of the parking lot management system.
type ParkingService struct {
	DB *gorm.DB
}

// NewParkingService creates and returns a new instance of ParkingService.
// It takes a pointer to a gorm.DB as a parameter to initialize the database connection.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance for database operations.
//
// Returns:
//   - *ParkingService: A new instance of ParkingService.
func NewParkingService(db *gorm.DB) *ParkingService {
	return &ParkingService{DB: db}
}

// CalculateParkingFee calculates the parking fee for a given parking entry.
// It retrieves the parking entry and associated tariff, then calculates the fee
// based on the duration of stay and the applicable rates.
//
// Parameters:
//   - entryID: The unique identifier of the parking entry.
//
// Returns:
//   - float64: The calculated parking fee.
//   - error: An error if any issues occur during the calculation process.
func (s *ParkingService) CalculateParkingFee(entryID uint) (float64, error) {
	var entry models.ParkingEntry
	if err := s.DB.First(&entry, entryID).Error; err != nil {
		return 0, err
	}

	if entry.ExitTime == nil {
		return 0, errors.ErrVehicleNotExited
	}

	var tariff models.Tariff
	if err := s.DB.Where("parking_lot_id = ? AND vehicle_type = ?", entry.ParkingLotID, entry.VehicleType).First(&tariff).Error; err != nil {
		return 0, errors.ErrTariffNotFound
	}

	duration := entry.ExitTime.Sub(entry.EntryTime)
	hours := int(math.Ceil(duration.Hours()))

	var fee float64

	if hours <= tariff.BaseHours {
		fee = tariff.BaseRate
	} else if hours <= tariff.DailyRateHours {
		fee = tariff.BaseRate + float64(hours-tariff.BaseHours)*tariff.HourlyRate
	} else {
		days := int(math.Ceil(float64(hours) / 24))
		fee = float64(days) * tariff.DailyRate
	}

	return fee, nil
}

// GetAvailableSpots returns the number of available spots for each vehicle type in a parking lot.
// It calculates the difference between the total capacity and the number of occupied spots.
//
// Parameters:
//   - parkingLotID: The unique identifier of the parking lot.
//
// Returns:
//   - map[models.VehicleType]int: A map containing the number of available spots for each vehicle type.
//   - error: An error if any issues occur during the retrieval process.
func (s *ParkingService) GetAvailableSpots(parkingLotID uint) (map[models.VehicleType]int, error) {
	var parkingLot models.ParkingLot
	if err := s.DB.First(&parkingLot, parkingLotID).Error; err != nil {
		return nil, errors.ErrParkingLotNotFound
	}

	availableSpots := make(map[models.VehicleType]int)
	for vehicleType, capacity := range parkingLot.Capacity {
		var count int64
		s.DB.Model(&models.ParkingEntry{}).Where("parking_lot_id = ? AND vehicle_type = ? AND exit_time IS NULL", parkingLotID, vehicleType).Count(&count)
		availableSpots[vehicleType] = capacity - int(count)
	}

	return availableSpots, nil
}

// CreateParkingEntry creates a new parking entry after checking for available spots.
// It ensures that there is space available for the given vehicle type before creating the entry.
//
// Parameters:
//   - entry: A pointer to the ParkingEntry struct containing the entry details.
//
// Returns:
//   - error: An error if the entry cannot be created or if there are no available spots.
func (s *ParkingService) CreateParkingEntry(entry *models.ParkingEntry) error {
	availableSpots, err := s.GetAvailableSpots(entry.ParkingLotID)
	if err != nil {
		return err
	}

	if availableSpots[entry.VehicleType] <= 0 {
		return errors.ErrNoAvailableSpots
	}

	entry.EntryTime = time.Now()
	return s.DB.Create(entry).Error
}

// ExitParking processes a vehicle exiting the parking lot.
// It updates the exit time of the parking entry, calculates the fee, and generates a receipt.
//
// Parameters:
//   - entryID: The unique identifier of the parking entry to be processed for exit.
//
// Returns:
//   - *models.ParkingEntry: A pointer to the updated ParkingEntry.
//   - *models.Receipt: A pointer to the generated Receipt.
//   - error: An error if any issues occur during the exit process.
func (s *ParkingService) ExitParking(entryID uint) (*models.ParkingEntry, *models.Receipt, error) {
	var entry models.ParkingEntry
	if err := s.DB.First(&entry, entryID).Error; err != nil {
		return nil, nil, errors.ErrParkingEntryNotFound
	}

	if entry.ExitTime != nil {
		return nil, nil, errors.ErrVehicleAlreadyExited
	}

	now := time.Now()
	entry.ExitTime = &now

	if err := s.DB.Save(&entry).Error; err != nil {
		return nil, nil, err
	}

	fee, err := s.CalculateParkingFee(entryID)
	if err != nil {
		return nil, nil, err
	}

	receipt := &models.Receipt{
		ParkingEntryID: entry.ID,
		EntryTime:      entry.EntryTime,
		ExitTime:       *entry.ExitTime,
		Duration:       entry.ExitTime.Sub(entry.EntryTime).String(),
		Fee:            fee,
	}

	if err := s.DB.Create(receipt).Error; err != nil {
		return nil, nil, err
	}

	return &entry, receipt, nil
}

// CreateTariff creates a new tariff in the database.
// It's used to set up pricing structures for different vehicle types in a parking lot.
//
// Parameters:
//   - tariff: A pointer to the Tariff struct containing the tariff details.
//
// Returns:
//   - error: An error if the tariff cannot be created in the database.
func (s *ParkingService) CreateTariff(tariff *models.Tariff) error {
	return s.DB.Create(tariff).Error
}