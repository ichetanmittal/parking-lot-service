package services

import (
	"math"
	"time"

	"github.com/chetan/parking-lot-service/internal/errors"
	"github.com/chetan/parking-lot-service/internal/models"
	"gorm.io/gorm"
)

type ParkingService struct {
	DB *gorm.DB
}

func NewParkingService(db *gorm.DB) *ParkingService {
	return &ParkingService{DB: db}
}

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
		if err == gorm.ErrRecordNotFound {
			return 0, errors.ErrTariffNotFound
		}
		return 0, err
	}

	duration := entry.ExitTime.Sub(entry.EntryTime)
	hours := math.Ceil(duration.Hours()) // Round up to the nearest hour

	var fee float64

	switch entry.VehicleType {
	case models.MotorcycleScooter:
		fee = hours * tariff.HourlyRate
	case models.CarSUV:
		if tariff.BaseHours > 0 && tariff.BaseRate > 0 {
			// Parking Lot B style: fixed rate for first hour, then hourly rate
			if hours <= float64(tariff.BaseHours) {
				fee = tariff.BaseRate
			} else {
				fee = tariff.BaseRate + (hours - float64(tariff.BaseHours)) * tariff.HourlyRate
			}
		} else {
			// Parking Lot A style: simple hourly rate
			fee = hours * tariff.HourlyRate
		}
	case models.BusTruck:
		if tariff.DailyRateHours > 0 && tariff.DailyRate > 0 {
			// Parking Lot A style: hourly rate up to a day, then daily rate
			if hours <= float64(tariff.DailyRateHours) {
				fee = hours * tariff.HourlyRate
			} else {
				days := math.Ceil(hours / 24)
				fee = float64(tariff.DailyRateHours) * tariff.HourlyRate + (days - 1) * tariff.DailyRate
			}
		} else {
			// Parking Lot B style: simple hourly rate
			fee = hours * tariff.HourlyRate
		}
	}

	return math.Round(fee*100) / 100, nil // Round to 2 decimal places
}

func (s *ParkingService) GetAvailableSpots(parkingLotID uint) (map[models.VehicleType]int, error) {
	var parkingLot models.ParkingLot
	if err := s.DB.First(&parkingLot, parkingLotID).Error; err != nil {
		return nil, err
	}

	type Result struct {
		VehicleType models.VehicleType
		Count       int
	}
	var results []Result

	if err := s.DB.Model(&models.ParkingEntry{}).
		Select("vehicle_type, count(*) as count").
		Where("parking_lot_id = ? AND exit_time IS NULL", parkingLotID).
		Group("vehicle_type").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	occupiedSpots := make(map[models.VehicleType]int)
	for _, result := range results {
		occupiedSpots[result.VehicleType] = result.Count
	}

	availableSpots := make(map[models.VehicleType]int)
	for vType, capacity := range parkingLot.Capacity {
		availableSpots[vType] = capacity - occupiedSpots[vType]
	}

	return availableSpots, nil
}

func (s *ParkingService) CreateParkingEntry(entry *models.ParkingEntry) error {
	// Check if parking lot exists
	var parkingLot models.ParkingLot
	if err := s.DB.First(&parkingLot, entry.ParkingLotID).Error; err != nil {
		return errors.ErrParkingLotNotFound
	}

	// Check if there are available spots
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

func (s *ParkingService) ExitParking(entryID uint) (*models.ParkingEntry, *models.Receipt, error) {
	var entry models.ParkingEntry
	if err := s.DB.First(&entry, entryID).Error; err != nil {
		return nil, nil, errors.ErrParkingEntryNotFound
	}

	if entry.ExitTime != nil {
		return nil, nil, errors.ErrVehicleAlreadyExited
	}

	exitTime := time.Now()
	entry.ExitTime = &exitTime

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

	// Save the receipt to the database
	if err := s.DB.Create(receipt).Error; err != nil {
		return nil, nil, err
	}

	return &entry, receipt, nil
}

// Add this method to the ParkingService struct
func (s *ParkingService) CreateTariff(tariff *models.Tariff) error {
	// Implement the logic to create a tariff
	return s.DB.Create(tariff).Error
}