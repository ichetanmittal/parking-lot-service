# Parking Lot Management System

## Overview

This project is a Parking Lot Management System implemented in Go. It provides a RESTful API for managing parking lots, vehicle entries and exits, and calculating parking fees based on customizable tariffs.

## Features

- Create and manage parking lots with customizable capacities for different vehicle types
- Track vehicle entries and exits
- Calculate parking fees based on vehicle type and duration
- Generate receipts for parking sessions
- RESTful API for easy integration

## Technology Stack

- Go 1.23
- Echo framework for HTTP routing
- GORM for database operations
- PostgreSQL as the database
- Docker for containerization

## API Documentation

The API is documented using OpenAPI 3.0 specification. You can find the full API documentation in the `api/openapi.yaml` file.

## Getting Started

### Prerequisites

- Go 1.23 or later
- PostgreSQL
- Docker (optional)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=parking_system
DB_PORT=5432
PORT=8080
```

### Running the Application

1. Clone the repository
2. Set up the environment variables
3. Run the following commands:

```bash
go mod download
go run cmd/server/main.go
```

### Running with Docker

1. Build the Docker image:
   ```
   docker build -t parking-lot-service .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 --env-file .env parking-lot-service
   ```

## API Endpoints

- `POST /parking-lots`: Create a new parking lot
- `GET /parking-lots/:id`: Get parking lot details
- `GET /parking-lots/:id/available-spots`: Get available spots in a parking lot
- `POST /parking-entries`: Create a new parking entry
- `PUT /parking-entries/:id/exit`: Process a vehicle exit and generate a receipt
- `POST /tariffs`: Create a new tariff

For detailed API documentation, refer to the `api/openapi.yaml` file.

## Models

The system uses the following main models:

- `ParkingLot`: Represents a parking lot with capacity for different vehicle types
- `ParkingEntry`: Represents a vehicle's parking session
- `Tariff`: Defines the pricing structure for different vehicle types
- `Receipt`: Generated when a vehicle exits, containing fee details

For detailed model structures, refer to the `api/openapi.yaml` file.

## Error Handling

Custom errors are defined in the `internal/errors/errors.go` file to handle various scenarios such as invalid input, not found errors, and business logic errors.

## Database

The system uses PostgreSQL as the database. Connection and migrations are handled in the `internal/db/db.go` file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
