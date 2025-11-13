# Quick Start Guide

Get the application running in 5 minutes.

## Prerequisites

- Go 1.24.4+
- Docker and Docker Compose
- Make

## Steps

### 1. Clone and Setup

```bash
git clone https://github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal.git
cd golang-ddd-modular-monolith-with-hexagonal
cp .env.example .env
```

### 2. Install Tools

```bash
make install-deps
```

### 3. Start Database

```bash
make dev-env
```

Wait for PostgreSQL to be ready.

### 4. Run Migrations

Open a new terminal and run:

```bash
cd golang-ddd-modular-monolith-with-hexagonal
echo "" | make migrate-up
```

### 5. Start Application

In the same terminal (or another one):

```bash
make dev-air
```

### 6. Test It

```bash
curl http://localhost:9090/health
```

Expected:

```json
{ "status": "ok", "mode": "rest", "environment": "development" }
```

## Try the API

### Create a Payment Setting

```bash
curl -X POST http://localhost:9090/api/v1/payment-settings \
  -H "Content-Type: application/json" \
  -d '{
    "settingKey": "min_amount",
    "settingValue": "10.00",
    "currency": "USD",
    "status": "active",
    "createdAt": "2025-11-13T10:00:00Z",
    "updatedAt": "2025-11-13T10:00:00Z"
  }'
```

### Create a Payment

```bash
curl -X POST http://localhost:9090/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.50,
    "currency": "USD",
    "status": "pending",
    "createdAt": "2025-11-13T10:00:00Z",
    "updatedAt": "2025-11-13T10:00:00Z"
  }'
```

### List Payments

```bash
curl http://localhost:9090/api/v1/payments
```

## Stop Everything

```bash
make down
```

## Full Documentation

See [README.md](README.md) for complete documentation.

## Postman Collection

[Access Postman Workspace](https://www.postman.com/crimson-shadow-8849/workspace/golang-ddd-modular-monolith-with-hexagonal/request/451883-2c3349ff-d54b-474d-bd02-9e25fe2efc69?action=share&creator=451883)
