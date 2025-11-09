package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Get database configuration from environment variables with defaults
	dbHost := getEnv("POSTGRES_HOST", "127.0.0.1")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "user")
	dbPassword := getEnv("POSTGRES_PASSWORD", "password")
	dbName := getEnv("POSTGRES_DB", "payment")
	dbSSLMode := getEnv("POSTGRES_SSLMODE", "disable")

	// Create PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Initialize database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Initialize modules in dependency order (like Nest.js)
	// Factory handles ALL the wiring: hexagon core + adapters

	// 1. Payment Settings module (no dependencies)
	paymentSettingsModule := settingsfactory.NewModule(settingsfactory.ModuleConfig{
		DB: db,
	})

	// 2. Payment module (depends on Payment Settings)
	paymentModule := paymentfactory.NewModule(paymentfactory.ModuleConfig{
		DB:                  db,
		PaymentSettingsPort: paymentSettingsModule.Service,
	})

	// Setup HTTP server and register module routes
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register module HTTP handlers (inbound adapters)
	api := e.Group("/api/v1")
	paymentModule.RegisterHTTPHandlers(api)
	paymentSettingsModule.RegisterHTTPHandlers(api)

	// Start server
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
