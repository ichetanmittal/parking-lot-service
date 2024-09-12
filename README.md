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

## Postman Collection

A Postman collection is available for testing and interacting with the API. It includes all the endpoints and sample requests:

1. **Endpoints in the Collection**:
   - `POST /parking-lots`: Create a new parking lot
   - `GET /parking-lots/:id`: Get parking lot details
   - `GET /parking-lots/:id/available-spots`: Get available spots in a parking lot
   - `POST /parking-entries`: Create a new parking entry
   - `GET /parking-entries/:id/exit`: Process a vehicle exit and generate a receipt
   - `POST /tariffs`: Create a new tariff

2. **How to Use the Collection**:
   - Download the Postman collection JSON file from the [`postman-collections` directory](./postman-collections/Parking-lot.postman_collection.json) in this repository.
   - Open Postman and click on "Import" in the top left corner.
   - Select the downloaded JSON file to import the collection into your Postman workspace.

## Getting Started

### Prerequisites

- Go 1.23 or later
- PostgreSQL
- Docker (optional)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

DB_HOST=localhost
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=parking_system
DB_PORT=5432
PORT=8080

### Running the Application

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/parking-lot-service.git
   ```

2. Navigate to the project directory:
   ```
   cd parking-lot-service
   ```

3. Install dependencies:
   ```
   go mod download
   ```

4. Run the application:
   ```
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080` (or the port specified in your .env file).

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

For detailed request/response formats, please refer to the API documentation.

## Testing

Run the tests using the following command:

```
go test ./...
```

To run tests with coverage:

```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)