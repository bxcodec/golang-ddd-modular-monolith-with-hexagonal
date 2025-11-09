# Project Setup Guide

This DDD Modular Monolith project is now configured to use PostgreSQL with Squirrel query builder.

## Prerequisites

- Docker and Docker Compose
- Go 1.24.4 or later
- Make

## Quick Start

### 1. Start PostgreSQL Database

```bash
make up
```

This will start PostgreSQL in Docker with the following default configuration:
- Host: `127.0.0.1`
- Port: `5432`
- User: `user`
- Password: `password`
- Database: `payment`

### 2. Run Database Migrations

```bash
make migrate-up
```

When prompted, press Enter to run all migrations (or specify a number).

This will create the following tables:
- `payment_settings` - stores payment configuration settings
- `payments` - stores payment transactions

### 3. Run the Application

```bash
make dev-air
```

Or run directly:

```bash
go run application/main.go
```

The server will start on `http://localhost:8080`

## Database Configuration

The application reads database configuration from environment variables:
- `POSTGRES_HOST` (default: `127.0.0.1`)
- `POSTGRES_PORT` (default: `5432`)
- `POSTGRES_USER` (default: `user`)
- `POSTGRES_PASSWORD` (default: `password`)
- `POSTGRES_DB` (default: `payment`)
- `POSTGRES_SSLMODE` (default: `disable`)

## Available Make Commands

### Development
- `make up` - Start Docker Compose (postgres) and air
- `make down` - Stop Docker Compose
- `make destroy` - Teardown (removes volumes, tmp files, etc.)
- `make dev-env` - Start only PostgreSQL container
- `make dev-air` - Start air (hot reload)

### Database Migrations
- `make migrate-up` - Apply migrations (interactive)
- `make migrate-down` - Rollback migrations (interactive)
- `make migrate-create` - Create new migration file
- `make migrate-drop` - Drop all tables (use with caution!)

### Building & Testing
- `make build` - Build the application binary
- `make build-race` - Build with race detector
- `make tests` - Run tests
- `make lint` - Run linter

### Tools Installation
- `make install-deps` - Install all development dependencies locally

## API Endpoints

### Payment Settings
- `GET /api/v1/payment-settings` - List payment settings
- `POST /api/v1/payment-settings` - Create payment setting
- `PUT /api/v1/payment-settings/:id` - Update payment setting
- `DELETE /api/v1/payment-settings/:id` - Delete payment setting

### Payments
- `GET /api/v1/payments` - List payments
- `GET /api/v1/payments/:id` - Get payment by ID
- `POST /api/v1/payments` - Create payment
- `PUT /api/v1/payments/:id` - Update payment
- `DELETE /api/v1/payments/:id` - Delete payment

## Project Structure

```
.
├── application/           # Application entry point
│   └── main.go           # Main application file with DB setup
├── modules/              # Domain modules
│   ├── payment/          # Payment module
│   │   ├── internal/
│   │   │   ├── adapter/
│   │   │   │   ├── controller/  # HTTP handlers
│   │   │   │   └── repository/  # Database layer (using Squirrel)
│   │   │   ├── ports/           # Interfaces
│   │   │   └── service/         # Business logic
│   │   └── factory/             # Module factory
│   └── payment-settings/ # Payment Settings module (similar structure)
├── migrations/           # SQL migration files
├── docker-compose.yml    # Docker services configuration
└── Makefile             # Build and development commands
```

## Implementation Details

### Query Builder
The project uses [Squirrel](https://github.com/Masterminds/squirrel) for building SQL queries:
- Fluent API for constructing SQL queries
- PostgreSQL placeholder format (`$1`, `$2`, etc.)
- Type-safe query building
- Direct query execution with `RunWith(db)`

### Repository Pattern
Both repositories (`PaymentRepository` and `PaymentSettingsRepository`) implement:
- CRUD operations using Squirrel query builder
- Cursor-based pagination for list operations
- Proper error handling
- Transaction support through `*sql.DB`

### Database Schema
Tables include:
- `id` (VARCHAR) - Primary key
- `amount` (DECIMAL) - Payment amount
- `currency` (VARCHAR) - Currency code (e.g., USD, EUR)
- `status` (VARCHAR) - Status indicator
- `created_at` (TIMESTAMP) - Creation timestamp
- `updated_at` (TIMESTAMP) - Last update timestamp

Indexes are created on frequently queried fields (currency, status, created_at).

## Troubleshooting

### Connection Issues
If you can't connect to the database:
1. Ensure Docker is running: `docker ps`
2. Check PostgreSQL logs: `docker-compose logs postgres`
3. Verify the database is healthy: `docker-compose ps`

### Migration Issues
If migrations fail:
1. Check the database connection
2. Verify migration files in `./migrations/`
3. Check current migration version: `make migrate-version` (if added)

### Port Conflicts
If port 5432 is already in use:
1. Stop the conflicting service
2. Or modify `docker-compose.yml` to use a different port

## Next Steps

1. Implement business logic in service layers
2. Add validation in controllers
3. Add unit tests for repositories and services
4. Add integration tests
5. Implement authentication/authorization
6. Add API documentation (Swagger/OpenAPI)

