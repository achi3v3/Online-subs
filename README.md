# DEMO Online Subscriptions Service

## Overview

RESTful service for aggregating data about users subscriptions. CRUDL-operations.

## Tech Stack

- **Language**: Go (Golang) 1.24+
- **Framework**: Gin Web Framework (github.com/gin-gonic/gin)
- **Database**: PostgreSQL (gorm.io/driver/postgres)
- **Database ORM**: GORM (gorm.io/gorm)
- **API Documentation**: Swagger UI (github.com/swaggo/gin-swagger)
- **Configuration**: godotenv for .env file support
- **Logging**: logrus (github.com/sirupsen/logrus)
- **Containerization**: Docker & Docker Compose

## API

View API documentation at `http://localhost:8080/swagger/index.html`

### Endpoints

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
2. Configure environment variables in `.env`  | `docker-compose.yml` files
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
