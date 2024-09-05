package errors

import "errors"

var (
	// ErrInvalidInput is returned when the input data is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrParkingLotNotFound is returned when a parking lot is not found
	ErrParkingLotNotFound = errors.New("parking lot not found")

	// ErrParkingEntryNotFound is returned when a parking entry is not found
	ErrParkingEntryNotFound = errors.New("parking entry not found")

	// ErrVehicleNotExited is returned when trying to calculate fees for a vehicle that hasn't exited
	ErrVehicleNotExited = errors.New("vehicle has not exited yet")

	// ErrTariffNotFound is returned when a tariff is not found for a specific parking lot and vehicle type
	ErrTariffNotFound = errors.New("tariff not found")

	// ErrNoAvailableSpots is returned when there are no available spots for a specific vehicle type
	ErrNoAvailableSpots = errors.New("no available spots for this vehicle type")

	// ErrVehicleAlreadyExited is returned when trying to process an exit for a vehicle that has already exited
	ErrVehicleAlreadyExited = errors.New("vehicle has already exited")
)
