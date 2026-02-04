# Go Microservice Template

A production-ready Go microservice template following Clean Architecture principles. Use this template to bootstrap new Go microservices with best practices already in place.

## Features

- ✅ **Clean Architecture** - Domain, UseCase, Infrastructure, and API layers
- ✅ **Generic CRUD** - Base repository, usecase, and handler with generics
- ✅ **Fiber Framework** - High-performance HTTP framework
- ✅ **GORM ORM** - Database operations with PostgreSQL
- ✅ **Swagger Docs** - Auto-generated API documentation
- ✅ **Prometheus Metrics** - Built-in observability
- ✅ **Structured Logging** - Zap logger with JSON output
- ✅ **Authentication** - Casdoor integration ready
- ✅ **Docker Ready** - Dockerfile and docker-compose included
- ✅ **CI/CD** - GitHub Actions workflow
- ✅ **Hot Reload** - Air configuration for development

## Quick Start

### 1. Create New Repository from Template

Click "Use this template" on GitHub or clone and reinitialize:

```bash
git clone https://github.com/minisource/template_go.git my-service
cd my-service
rm -rf .git
git init
```

### 2. Initialize Your Project

```bash
# Using Make (recommended)
make init MODULE=github.com/your-org/my-service SERVICE=my-service

# Or using the init script
./scripts/init.sh github.com/your-org/my-service my-service
```

### 3. Configure Environment

```bash
cp src/.env.example src/.env
# Edit src/.env with your database and auth settings
```

### 4. Run the Service

```bash
# Using Make
make run

# Or using Task
cd src && task run

# Or directly
cd src && go run ./cmd/main.go
```

## Project Structure

```
├── .github/
│   └── workflows/          # GitHub Actions CI/CD
├── docker/
│   ├── docker-compose.yml  # Full stack with ELK, Prometheus, Grafana
│   ├── alertmanager/       # Alert configuration
│   ├── elk/                # Elasticsearch, Kibana, Filebeat
│   ├── grafana/            # Dashboards and datasources
│   ├── prometheus/         # Metrics and alerts
│   └── redis/              # Redis configuration
├── docs/                   # Project documentation
├── src/
│   ├── api/
│   │   ├── api.go          # Server initialization
│   │   ├── dto/            # Data Transfer Objects (request/response)
│   │   ├── handler/        # HTTP handlers
│   │   ├── middleware/     # Custom middleware
│   │   ├── router/         # Route definitions
│   │   └── validation/     # Custom validators
│   ├── cmd/
│   │   └── main.go         # Application entry point
│   ├── config/
│   │   ├── config.go       # Configuration struct and loading
│   │   └── config-*.yml    # Environment-specific configs
│   ├── constant/           # Application constants
│   ├── dependency/         # Dependency injection setup
│   ├── domain/
│   │   ├── model/          # Domain entities
│   │   └── repository/     # Repository interfaces
│   ├── infra/
│   │   └── persistence/
│   │       ├── database/   # Database connection
│   │       ├── migration/  # Database migrations
│   │       └── repository/ # Repository implementations
│   ├── pkg/                # Shared packages (can be imported by other projects)
│   ├── tests/
│   │   ├── integration/    # Integration tests
│   │   └── unit/           # Unit tests
│   └── usecase/            # Business logic layer
├── Makefile                # Build and development commands
└── README.md
```

## Architecture

This template follows **Clean Architecture** (also known as Hexagonal Architecture):

```
┌─────────────────────────────────────────────────────────────┐
│                         API Layer                           │
│  (Handlers, DTOs, Routers, Middleware, Validation)         │
├─────────────────────────────────────────────────────────────┤
│                       UseCase Layer                         │
│           (Business Logic, Orchestration)                   │
├─────────────────────────────────────────────────────────────┤
│                       Domain Layer                          │
│        (Entities, Repository Interfaces)                    │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                      │
│    (Database, External Services, Repository Impl)          │
└─────────────────────────────────────────────────────────────┘
```

**Key Principles:**
- Dependencies point inward (outer layers depend on inner layers)
- Domain layer has no external dependencies
- Business logic lives in UseCase layer, not handlers
- Infrastructure details are abstracted via interfaces

## Adding a New Entity

### 1. Create Domain Model

```go
// src/domain/model/product.go
package model

type Product struct {
    BaseModel
    Name        string  `gorm:"size:100;not null"`
    Description string  `gorm:"size:500"`
    Price       float64 `gorm:"not null"`
}
```

### 2. Create Repository Interface

```go
// src/domain/repository/product_repository.go
package repository

import "your-module/domain/model"

type ProductRepository interface {
    BaseRepository[model.Product]
    // Add custom methods here
    FindByName(ctx context.Context, name string) (*model.Product, error)
}
```

### 3. Create DTOs

```go
// src/api/dto/product.go
package dto

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"required,gt=0"`
}

type ProductResponse struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
```

### 4. Create UseCase

```go
// src/usecase/product_usecase.go
package usecase

import (
    "your-module/config"
    "your-module/domain/repository"
)

type ProductUsecase struct {
    *BaseUsecase[model.Product, dto.CreateProductRequest, dto.UpdateProductRequest, dto.ProductResponse]
}

func NewProductUsecase(cfg *config.Config, repo repository.ProductRepository) *ProductUsecase {
    return &ProductUsecase{
        BaseUsecase: NewBaseUsecase[...](cfg, repo),
    }
}
```

### 5. Create Handler and Router

```go
// src/api/handler/product.go
// src/api/router/product.go
```

### 6. Wire Dependencies

```go
// src/dependency/dependency.go
func GetProductRepository(cfg *config.Config) repository.ProductRepository {
    return infrarepository.NewBaseRepository[model.Product](cfg, nil)
}
```

### 7. Add Migration

```go
// src/infra/persistence/migration/2_AddProduct.go
tables = addNewTable(database, model.Product{}, tables)
```

## Available Commands

```bash
make help              # Show all available commands
make init              # Initialize new project
make build             # Build the application
make run               # Run the application
make test              # Run all tests
make test-unit         # Run unit tests only
make test-integration  # Run integration tests
make lint              # Run linter
make fmt               # Format code
make swagger           # Generate Swagger docs
make docker-build      # Build Docker image
make docker-run        # Run with Docker Compose
make install-tools     # Install dev tools (golangci-lint, swag, air)
```

## Configuration

Configuration is loaded from YAML files based on `APP_ENV`:

| APP_ENV     | Config File              |
|-------------|--------------------------|
| development | config-development.yml   |
| docker      | config-docker.yml        |
| production  | config-production.yml    |

Environment variables can override config values.

## API Documentation

Swagger UI is available at: `http://localhost:5005/swagger/`

To regenerate docs after API changes:
```bash
make swagger
```

## Monitoring

The template includes full observability stack:

- **Prometheus**: `http://localhost:9090` - Metrics
- **Grafana**: `http://localhost:3000` - Dashboards
- **Kibana**: `http://localhost:5601` - Logs
- **App Metrics**: `http://localhost:5005/metrics`

## Testing

```bash
# Run all tests with coverage
make test

# Run specific test
cd src && go test -v ./usecase/... -run TestProductCreate
```

## Docker

```bash
# Build image
make docker-build

# Run full stack (app + postgres + redis + monitoring)
make docker-run

# Stop all containers
make docker-stop
```

## Dependencies

Core dependencies managed via `go-common`:
- **Fiber** - HTTP framework
- **GORM** - ORM
- **Viper** - Configuration
- **Zap** - Logging
- **Prometheus** - Metrics

## License

MIT License - see [LICENSE](LICENSE) file