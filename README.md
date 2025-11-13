# Go DDD Modular Monolith with Hexagonal Architecture

[![Go Version](https://img.shields.io/badge/Go-1.24.4+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Golang example project demonstrating Domain-Driven Design (DDD) in a Modular Monolith architecture with Hexagonal (Ports and Adapters) pattern.

## Read the original Blogpost here:

- [Combining Modular Monolith and Hexagonal Architecture while Maintaining Domain Driven Design Principles (part 1)](https://notes.softwarearchitect.id/p/combining-modular-monolith-and-hexagonal)
- [Developing Modular Monolith and Hexagonal Architecture in Golang while Maintaining Domain Driven Design Principles (part 2)](https://notes.softwarearchitect.id/p/developing-modular-monolith-and-hexagonal)

Please read the first article, then proceed to the second and this repository to understand the reasoning behind the current structure.

**Quick Links:**

- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Postman Collection](https://www.postman.com/crimson-shadow-8849/workspace/golang-ddd-modular-monolith-with-hexagonal/request/451883-2c3349ff-d54b-474d-bd02-9e25fe2efc69?action=share&creator=451883) - API documentation

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [CLI Commands](#cli-commands)
- [Development](#development)
- [Database Migrations](#database-migrations)
- [Docker](#docker)
- [Testing](#testing)
- [Architecture Benefits](#architecture-benefits)
- [Production Deployment](#production-deployment)
- [Contributing](#contributing)
- [License](#license)

## Architecture Overview

This project combines three powerful architectural patterns:

### Modular Monolith

- Multiple independent modules (payment, payment-settings) in a single deployable unit
- Each module has clear boundaries and communicates via well-defined interfaces
- Modules can be developed, tested, and evolved independently
- All modules share the same process and database, simplifying deployment and transactions

### Hexagonal Architecture (Ports & Adapters)

- Domain logic (hexagon core) is isolated from external concerns
- Inbound adapters: REST controllers, cron jobs (drive the application)
- Outbound adapters: Repositories, external APIs (driven by the application)
- Ports: Interfaces that define contracts between core and adapters

### Domain-Driven Design (DDD)

- Clear separation between domain logic and infrastructure
- Each module represents a bounded context
- Domain entities and business rules are at the center

## Project Structure

```
.
├── application/            # Application entry point
├── cmd/                    # CLI commands (Cobra)
│   ├── rest.go            # REST API server command
│   ├── cron_update_payment.go  # Cron job command
│   └── root.go            # Root command configuration
├── modules/               # Business modules (bounded contexts)
│   ├── payment/
│   │   ├── factory/       # Module factory for dependency injection
│   │   ├── internal/
│   │   │   ├── adapter/   # Inbound/Outbound adapters
│   │   │   │   ├── controller/  # REST controllers (inbound)
│   │   │   │   ├── cron/        # Cron jobs (inbound)
│   │   │   │   └── repository/  # Database repository (outbound)
│   │   │   ├── ports/     # Interface definitions
│   │   │   └── service/   # Business logic
│   │   ├── module.go      # Module registration
│   │   └── payment.go     # Domain entities / Public API for the domain/module
│   └── payment-settings/  # Similar structure
├── pkg/                   # Shared utilities
│   ├── config/            # Configuration management
│   ├── dbutils/           # Database utilities
│   ├── errors/            # Error handling
│   ├── logger/            # Logging utilities
│   ├── middlewares/       # HTTP middlewares
│   └── uniqueid/          # ID generation (ULID)
├── migrations/            # Database migrations
└── docker-compose.yml     # Docker setup
```

## Features

- Clean Architecture with clear separation of concerns
- Hexagonal Architecture (Ports and Adapters)
- Domain-Driven Design principles
- RESTful API using Echo framework
- PostgreSQL database with migrations
- Structured logging with zerolog
- Comprehensive testing (unit and E2E)
- Docker support
- Hot reload development with Air
- CLI commands with Cobra
- Mock generation with Mockery
- Code linting with golangci-lint

## Prerequisites

- Go 1.24.4 or higher
- Docker and Docker Compose (for database)
- Make (for running commands)

## Quick Start

### Option A: One-Command Setup (Recommended for First Time)

If you want everything set up automatically:

```bash
# Clone the repository
git clone https://github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal.git
cd golang-ddd-modular-monolith-with-hexagonal

# Copy environment file
cp .env.example .env

# Install dependencies and start everything
make init
make up
```

The `make up` command will:

1. Start PostgreSQL in Docker
2. Start the application with hot reload (you'll need to run migrations separately, see below)

After the app starts, open a new terminal and run migrations:

```bash
# In a new terminal window
cd golang-ddd-modular-monolith-with-hexagonal
make migrate-up
# Press Enter when prompted to apply all migrations
```

The API will be available at `http://localhost:9090`

```bash
curl http://localhost:9090/health
```

Expected response:

```json
{
  "status": "ok",
  "mode": "rest",
  "environment": "development"
}
```

## CLI Commands

The application supports multiple commands via Cobra:

### Start REST API Server

```bash
go run application/main.go rest
```

### Run Cron Job

```bash
go run application/main.go cron-update-payment
```

Options:

- `--batch-size` - Number of payments to process in one batch
- `--dry-run` - Run in dry-run mode without making changes

Example:

```bash
go run application/main.go cron-update-payment --batch-size 100 --dry-run
```

## Development

### Hot Reload with Air

Air is automatically started when you run `make up`. It watches for file changes and rebuilds the application.

### Running Tests

```bash
# Run unit tests
make test-unit

# Run E2E tests
make test-e2e

# Run all tests
make test-all

# Run tests with detailed output
make tests-complete
```

### Code Formatting

```bash
make fmt
```

This will format your code using:

- gofmt
- gofumpt
- goimports

### Linting

```bash
make lint
```

### Generate Mocks

```bash
make go-generate
```

This uses Mockery to generate mocks based on the `.mockery.yml` configuration.

## Database Migrations

### Apply Migrations

```bash
make migrate-up
```

### Rollback Migrations

```bash
make migrate-down
```

### Create New Migration

```bash
make migrate-create
```

### Drop Database

```bash
make migrate-drop
```

## Docker

### Build Docker Image

```bash
make image-build
```

### Run with Docker Compose

```bash
docker compose up
```

This will start both PostgreSQL and the application.

## Makefile Commands

Run `make help` to see all available commands:

```bash
make help
```

Key commands:

- `make up` - Start development environment
- `make down` - Stop Docker containers
- `make destroy` - Teardown everything (removes volumes)
- `make build` - Build the binary
- `make tests` - Run tests
- `make lint` - Run linter
- `make fmt` - Format code
- `make clean` - Clean artifacts

## Testing

The project includes both unit tests and E2E tests:

### Unit Tests

Located alongside the code, testing individual components in isolation using mocks.

### E2E Tests

Use testcontainers to spin up real PostgreSQL instances for integration testing.

### Test Coverage

```bash
make test-all
```

Coverage reports are generated in `coverage-all.out`.

## Architecture Benefits

### Why Modular Monolith?

- Simpler than microservices while maintaining modularity
- Easy refactoring to microservices if needed (modules are already isolated)
- Shared infrastructure reduces operational complexity
- No network latency between modules

### Why Hexagonal Architecture?

- Business logic is independent of frameworks and tools
- Easy to swap implementations (e.g., change database, HTTP framework)
- Highly testable (mock adapters, test ports)
- Clear separation between what the application does and how it does it

### Why DDD?

- Focus on core domain and domain logic
- Ubiquitous language between developers and domain experts
- Clear bounded contexts
- Better code organization and maintainability

## Module Communication

Modules communicate via well-defined ports (interfaces):

```go
// payment module depends on payment-settings module
type IPaymentSettingsPort interface {
    GetPaymentSetting(id string) (PaymentSetting, error)
    // ...
}

// payment module uses the interface, not the concrete implementation
paymentModule := paymentfactory.NewModule(paymentfactory.ModuleConfig{
    DB:                  db,
    PaymentSettingsPort: paymentSettingsModule.Service, // dependency injection
})
```

This allows:

- Clear contracts between modules
- Easy testing with mocks
- Potential extraction to microservices

## Production Deployment

### Build Binary

```bash
make build
```

The binary will be created as `engine` in the project root.

### Run Binary

```bash
./engine rest
```

### Docker Deployment

```bash
docker build -t payment-app:latest .
docker run -p 9090:9090 \
  -e POSTGRES_HOST=your-db-host \
  -e POSTGRES_PASSWORD=your-db-password \
  payment-app:latest rest
```

## Contributing

Contributions are welcome. Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Write tests for new features
4. Ensure all tests pass
5. Run linter and formatter
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Resources

- Clean Architecture: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- Hexagonal Architecture: https://alistair.cockburn.us/hexagonal-architecture/
- Hexagonal Architecture Example: https://github.com/jmgarridopaz/bluezone
- Presentation Domain Data Layering: https://martinfowler.com/bliki/PresentationDomainDataLayering.html
- Original Clean Architecture Example: https://github.com/bxcodec/go-clean-arch
- Modular Monolith in .NET: https://github.com/kgrzybek/modular-monolith-with-ddd
- Architecting Robust .NET: https://medium.com/@mail2mhossain/architecting-robust-net-dfa4f3725142

## Author

Iman Tumorang

## Acknowledgments

This project is inspired by various clean architecture implementations and aims to demonstrate a practical approach to building maintainable Go applications.
