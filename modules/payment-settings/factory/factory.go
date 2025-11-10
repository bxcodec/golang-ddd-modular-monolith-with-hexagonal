// Package factory provides module initialization and dependency wiring.
//
// The factory pattern is used to assemble modules in a modular monolith:
//   - Instantiates all adapters (repositories, controllers)
//   - Wires dependencies through constructor injection
//   - Returns a complete, ready-to-use Module
//
// This approach keeps the wiring logic separate from the domain and allows
// different configurations for different environments (dev, test, prod).
package factory

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/controller"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/repository"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/service"
)

// ModuleConfig contains all external dependencies required to initialize the Payment Settings module.
type ModuleConfig struct {
	DB *sql.DB
}

// NewModule assembles and wires the complete Payment Settings module using dependency injection.
//
// This is where hexagonal architecture comes together:
//   1. Create outbound adapters (repository for database access)
//   2. Inject adapters into the core service (hexagon)
//   3. Create inbound adapters (HTTP controller)
//   4. Return the module with all components connected
//
// The result is a fully independent module that can be deployed as part of a monolith.
func NewModule(config ModuleConfig) *paymentsettings.Module {
	// Wire up outbound adapters (repositories)
	settingsRepo := repository.NewPaymentSettingsRepository(config.DB)

	// Wire up the hexagon core (service)
	settingsService := service.NewPaymentSettingsService(settingsRepo)

	return &paymentsettings.Module{
		Service: settingsService,
		RegisterController: func(e *echo.Group) {
			controller.NewPaymentSettingController(e, settingsService)
		},
	}
}
