package factory

import (
	"database/sql"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/adapter/repository"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/service"
)

// paymentSettingsAdapter adapts the public payment settings service to the internal port interface
type paymentSettingsAdapter struct {
	service paymentsettings.IPaymentSettingsService
}

func (a *paymentSettingsAdapter) GetPaymentSettingsByCurrency(currency string) (paymentsettings.PaymentSettings, error) {
	return a.service.GetPaymentSettingsByCurrency(currency)
}

// NewPaymentService initializes and returns a payment service with all its dependencies wired up.
// This is the factory function that wires up all internal dependencies.
func NewPaymentService(db *sql.DB, paymentSettingsService paymentsettings.IPaymentSettingsService) payment.IPaymentService {
	paymentRepo := repository.NewPaymentRepository(db)
	settingsPort := &paymentSettingsAdapter{service: paymentSettingsService}
	return service.NewPaymentService(paymentRepo, settingsPort)
}
