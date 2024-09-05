package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// VehicleType represents the type of vehicle that can be parked.
type VehicleType string

const (
	MotorcycleScooter VehicleType = "MotorcycleScooter"
	CarSUV            VehicleType = "CarSUV"
	BusTruck          VehicleType = "BusTruck"
)

// ParkingLot represents a parking facility with its capacity for different vehicle types.
type ParkingLot struct {
	gorm.Model
	Name     string                `gorm:"not null"`
	Capacity map[VehicleType]int   `gorm:"serializer:json;not null"`
	Tariffs  []Tariff              `gorm:"foreignKey:ParkingLotID"`
}

// Tariff defines the pricing structure for a specific vehicle type in a parking lot.
type Tariff struct {
	gorm.Model
	ParkingLotID   uint        `gorm:"not null"`
	VehicleType    VehicleType `gorm:"not null"`
	BaseRate       float64     `gorm:"not null"`
	BaseHours      int         `gorm:"not null"`
	HourlyRate     float64     `gorm:"not null"`
	DailyRate      float64
	DailyRateHours int
}

// ParkingEntry represents a vehicle's entry into a parking lot.
type ParkingEntry struct {
	gorm.Model
	ParkingLotID uint        `gorm:"not null"`
	VehicleType  VehicleType `gorm:"not null"`
	LicensePlate string      `gorm:"not null"`
	EntryTime    time.Time   `gorm:"not null"`
	ExitTime     *time.Time
	Receipt      *Receipt    `gorm:"foreignKey:ParkingEntryID"`
}

// Receipt represents a parking fee receipt generated when a vehicle exits.
type Receipt struct {
	gorm.Model
	ParkingEntryID uint      `gorm:"not null"`
	EntryTime      time.Time `gorm:"not null"`
	ExitTime       time.Time `gorm:"not null"`
	Duration       string    `gorm:"not null"`
	Fee            float64   `gorm:"not null"`
}

// Validate checks if the ParkingLot data is valid.
func (pl *ParkingLot) Validate() error {
	if pl.Name == "" {
		return errors.New("invalid input: name is empty")
	}
	if len(pl.Capacity) == 0 {
		return errors.New("invalid input: capacity is empty")
	}
	for vType, capacity := range pl.Capacity {
		if capacity < 0 {
			return errors.New("invalid input: negative capacity for " + string(vType))
		}
	}
	return nil
}

// Validate checks if the Tariff data is valid.
func (t *Tariff) Validate() error {
	if t.ParkingLotID == 0 {
		return errors.New("invalid tariff input: missing parking lot ID")
	}
	if t.BaseRate < 0 || t.HourlyRate < 0 || t.DailyRate < 0 || t.BaseHours < 0 || t.DailyRateHours < 0 {
		return errors.New("invalid tariff input: negative rates or hours")
	}
	return nil
}

// Validate checks if the ParkingEntry data is valid.
func (pe *ParkingEntry) Validate() error {
	if pe.ParkingLotID == 0 {
		return errors.New("invalid parking entry input: missing parking lot ID")
	}
	if pe.LicensePlate == "" {
		return errors.New("invalid parking entry input: missing license plate")
	}
	if pe.ExitTime != nil && pe.ExitTime.Before(pe.EntryTime) {
		return errors.New("invalid parking entry input: exit time before entry time")
	}
	return nil
}

// Validate checks if the Receipt data is valid.
func (r *Receipt) Validate() error {
	if r.ParkingEntryID == 0 {
		return errors.New("invalid receipt input: missing parking entry ID")
	}
	if r.ExitTime.Before(r.EntryTime) {
		return errors.New("invalid receipt input: exit time before entry time")
	}
	if r.Fee < 0 {
		return errors.New("invalid receipt input: negative fee")
	}
	return nil
}
