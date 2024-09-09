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