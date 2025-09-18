# DEMO Online Subscriptions Service

## Overview

This project implements a RESTful service for aggregating data about users' online subscriptions. It provides CRUDL operations for subscription records and allows calculating the total cost of subscriptions over a given period with filtering capabilities.

## Tech Stack

- **Language**: Go (Golang) 1.24+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Database ORM**: GORM
- **API Documentation**: Swagger UI (gin-swagger)
- **Configuration**: godotenv for .env file support
- **Logging**: logrus
- **Containerization**: Docker & Docker Compose

## Features

1. CRUDL operations for subscription records:
   - Create, Read, Update, Delete, and List subscriptions
   - Each subscription includes:
     - Service name
     - Monthly price in RUB
     - User ID (UUID format)
     - Start date (month and year)
     - Optional end date
2. Calculate total subscription cost for a given period with filters:
   - Filter by user ID
   - Filter by service name
3. PostgreSQL database with migration support
4. Comprehensive logging throughout the application
5. Configuration via environment variables (.env file)
6. Swagger API documentation
7. Docker Compose for easy deployment

## API Endpoints

- `POST /subs` - Create a new subscription
- `GET /subs/:id` - Get a subscription by ID
- `PUT /subs/:id` - Update a subscription by ID
- `DELETE /subs/:id` - Delete a subscription by ID
- `GET /subs` - List subscriptions with optional filtering
- `GET /subs/total` - Calculate total subscription cost for a period
- `GET /` - Health check endpoint
- `GET /swagger/*any` - API documentation

## Getting Started

1. Clone the repository
2. Configure environment variables in `.env` file
3. Run `docker-compose up` to start the service and database
4. Access the API at `http://localhost:8080`
5. View API documentation at `http://localhost:8080/swagger/index.html`

## Project Structure

```
backend/
├── main.go          # Application entry point
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── .env             # Configuration file
├── Dockerfile       # Docker configuration
├── docs/            # Swagger documentation files
├── internal/        # Application source code
│   ├── database/    # Database initialization
│   ├── models/      # Data models
│   └── subs/        # Subscriptions business logic (handlers, service, repository)
└── migrations/      # Database migration scripts
```