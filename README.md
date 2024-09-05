# Parking Lot Management System

## Overview

This project is a Parking Lot Management System implemented in Go. It provides a RESTful API for managing parking lots, vehicle entries and exits, and calculating parking fees based on customizable tariffs.

## Features

- Create and manage parking lots with customizable capacities for different vehicle types
- Track vehicle entries and exits
- Calculate parking fees based on vehicle type and duration
- Generate receipts for parking sessions
- RESTful API for easy integration

## Project Structure

```
parking-lot-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── db/
│   │   └── db.go
│   ├── errors/
│   │   └── errors.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── models/
│   │   └── models.go
│   └── services/
│       └── parking_service.go
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go 1.16 or later
- Docker (optional, for containerized deployment)
- PostgreSQL database

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/chetan/parking-lot-service.git
   cd parking-lot-service
   ```

2. Set up environment variables:
   Create a `.env` file in the root directory with the following content:
   ```
   DB_HOST=localhost
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=parking_system
   DB_PORT=5432
   PORT=8080
   ```

3. Install dependencies:
   ```
   go mod download
   ```

4. Build the application:
   ```
   go build -o parking-service ./cmd/server
   ```

5. Run the application:
   ```
   ./parking-service
   ```

## Docker Deployment

To run the application using Docker:

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
- `GET /receipts/:id`: Get receipt details

## Models

The system uses the following main models:

- `ParkingLot`: Represents a parking lot with capacity for different vehicle types
- `ParkingEntry`: Represents a vehicle's parking session
- `Tariff`: Defines the pricing structure for different vehicle types
- `Receipt`: Generated when a vehicle exits, containing fee details

For detailed model structures, refer to the `models.go` file.

## Error Handling

Custom errors are defined in the `errors.go` file to handle various scenarios such as invalid input, not found errors, and business logic errors.

## Database

The system uses PostgreSQL as the database. Connection and migrations are handled in the `db.go` file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
