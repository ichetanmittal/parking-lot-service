package errors

import "errors"

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrParkingLotNotFound  = errors.New("parking lot not found")
	ErrParkingEntryNotFound = errors.New("parking entry not found")
	ErrVehicleNotExited    = errors.New("vehicle has not exited yet")
	ErrTariffNotFound      = errors.New("tariff not found")
	ErrNoAvailableSpots    = errors.New("no available spots for this vehicle type")
	ErrVehicleAlreadyExited = errors.New("vehicle has already exited")
)
