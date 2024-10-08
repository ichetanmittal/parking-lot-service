package models

import (
	"errors"
	"time"
	"gorm.io/gorm"
)

// VehicleType represents the type of vehicle
type VehicleType string

const (
	MotorcycleScooter VehicleType = "MotorcycleScooter"
	CarSUV            VehicleType = "CarSUV"
	BusTruck          VehicleType = "BusTruck"
)

// ParkingLot represents a parking lot in the system
type ParkingLot struct {
	gorm.Model
	Name     string
	Capacity map[VehicleType]int `gorm:"serializer:json"`
}

// Tariff represents the pricing structure for a specific vehicle type in a parking lot
type Tariff struct {
	gorm.Model
	ParkingLotID     uint
	VehicleType      VehicleType
	BaseRate         float64
	BaseHours        int
	HourlyRate       float64
	DailyRate        float64
	DailyRateHours   int
}

// ParkingEntry represents a vehicle's parking session
type ParkingEntry struct {
	gorm.Model
	ParkingLotID uint
	VehicleType  VehicleType
	LicensePlate string
	EntryTime    time.Time
	ExitTime     *time.Time
}

// Validate checks if the ParkingLot data is valid
func (pl *ParkingLot) Validate() error {
	if pl.Name == "" {
		return errors.New("invalid input: name is empty")
	}
	if len(pl.Capacity) == 0 {
		return errors.New("invalid input: capacity is empty")
	}
	return nil
}

// Validate checks if the Tariff data is valid
func (t *Tariff) Validate() error {
	if t.ParkingLotID == 0 || t.BaseRate < 0 || t.HourlyRate < 0 || t.DailyRate < 0 || t.BaseHours < 0 || t.DailyRateHours < 0 {
		return errors.New("invalid tariff input")
	}
	return nil
}

// Validate checks if the ParkingEntry data is valid
func (pe *ParkingEntry) Validate() error {
	if pe.ParkingLotID == 0 || pe.LicensePlate == "" {
		return errors.New("invalid parking entry input")
	}
	return nil
}

// Receipt represents a parking receipt generated when a vehicle exits
type Receipt struct {
	gorm.Model
	ParkingEntryID uint
	EntryTime      time.Time
	ExitTime       time.Time
	Duration       string
	Fee            float64
}