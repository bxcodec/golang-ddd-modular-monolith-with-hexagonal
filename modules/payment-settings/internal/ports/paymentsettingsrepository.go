package ports

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type IPaymentSettingsRepository interface {
	GetPaymentSettingsByCurrency(currency string) (paymentsettings.PaymentSettings, error)
	CreatePaymentSettings(settings *paymentsettings.PaymentSettings) error
	UpdatePaymentSettings(settings *paymentsettings.PaymentSettings) error
	DeletePaymentSettings(id string) error
}
