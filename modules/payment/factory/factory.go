// Package factory provides module initialization and dependency wiring.
//
// The factory pattern is used to assemble modules in a modular monolith:
//   - Instantiates all adapters (repositories, controllers, cron jobs)
//   - Wires dependencies through constructor injection
//   - Returns a complete, ready-to-use Module
//
// This approach keeps the wiring logic separate from the domain and allows
// different configurations for different environments (dev, test, prod).
package factory

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/controller"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/cron"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/repository"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/service"
)

// ModuleConfig contains all external dependencies required to initialize the Payment module.
//
// PaymentSettingsPort demonstrates inter-module dependencies in a modular monolith:
//   - The Payment module depends on Payment Settings functionality
//   - Instead of direct module import, we depend on a port interface
//   - The calling code injects the actual implementation
//   - This prevents circular dependencies and maintains module boundaries
type ModuleConfig struct {
	DB                  *sql.DB
	PaymentSettingsPort ports.IPaymentSettingsPort
	CronBatchSize       int
	CronDryRun          bool
}

// NewModule assembles and wires the complete Payment module using dependency injection.
//
// This is where hexagonal architecture comes together:
//   1. Create outbound adapters (repository for database access)
//   2. Inject adapters into the core service (hexagon)
//   3. Create inbound adapters (HTTP controller, cron jobs)
//   4. Return the module with all components connected
//
// The result is a fully independent module that can be deployed as part of a monolith
// or potentially extracted into a microservice with minimal changes.
func NewModule(config ModuleConfig) *payment.Module {
	// Wire up outbound adapters (repositories)
	paymentRepo := repository.NewPaymentRepository(config.DB)

	// Wire up the hexagon core (service)
	paymentService := service.NewPaymentService(paymentRepo, config.PaymentSettingsPort)

	// Set default cron batch size if not provided
	if config.CronBatchSize == 0 {
		config.CronBatchSize = 50
	}

	// Wire up cron adapters
	paymentUpdater := cron.NewPaymentUpdater(paymentService, cron.PaymentUpdaterConfig{
		BatchSize: config.CronBatchSize,
		DryRun:    config.CronDryRun,
	})

	// Create the module with all adapters
	return &payment.Module{
		Service: paymentService,
		RegisterController: func(e *echo.Group) {
			controller.NewPaymentController(e, paymentService)
		},
		PaymentUpdater: paymentUpdater,
	}
}
