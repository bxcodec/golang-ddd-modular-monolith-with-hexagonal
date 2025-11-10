package factory

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/controller"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/adapter/repository"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/service"
)

// ModuleConfig contains the dependencies needed to initialize the payment-settings module
type ModuleConfig struct {
	DB *sql.DB
}

// NewModule creates a fully wired payment-settings module
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
