package main

import (
	"database/sql"
	"log"

	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize modules in dependency order using factory packages
	// 1. Payment Settings module (no dependencies)
	paymentSettingsService := settingsfactory.NewPaymentSettingsService(db)

	// 2. Payment module (depends on Payment Settings)
	paymentService := paymentfactory.NewPaymentService(db, paymentSettingsService)

	// Use the services
	_ = paymentService
	_ = paymentSettingsService

	log.Println("Application initialized successfully!")
}
