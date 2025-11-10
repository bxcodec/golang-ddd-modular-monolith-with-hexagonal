package factory

import (
	"database/sql"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/controller"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/cron"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/repository"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/service"
	"github.com/labstack/echo/v4"
)

// ModuleConfig contains the dependencies needed to initialize the payment module
type ModuleConfig struct {
	DB *sql.DB
	// Payment Module depends on Payment Settings "public API".
	// Here we define the public API as the dependency needed to initialize the module.
	// But we don't need to import directly the module here. Instead we defined the port interface in the ports package.
	// this is to follow Golang idiomatic way of doing things.
	// also this is to avoid circular dependencies.
	PaymentSettingsPort ports.IPaymentSettingsPort
	// Cron configuration
	CronBatchSize int
	CronDryRun    bool
	// In the future if we need more dependencies, we can add them here.
}

// NewModule creates a fully wired payment module
// This assembles ALL components: core (hexagon) + adapters (controllers, repositories, cron)
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
