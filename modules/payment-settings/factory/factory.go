package factory

import (
	"database/sql"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/repository"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/service"
)

// NewPaymentSettingsService initializes and returns a payment settings service with all its dependencies wired up.
// This is the factory function that wires up all internal dependencies.
func NewPaymentSettingsService(db *sql.DB) paymentsettings.IPaymentSettingsService {
	settingsRepo := repository.NewPaymentSettingsRepository(db)
	return service.NewPaymentSettingsService(settingsRepo)
}
