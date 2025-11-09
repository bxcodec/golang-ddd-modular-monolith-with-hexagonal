package main

import (
	"database/sql"
	"log"

	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
