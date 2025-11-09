package ports

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

// IPaymentSettingsPort defines the contract for payment settings
type IPaymentSettingsPort interface {
	GetPaymentSettingsByCurrency(currency string) (paymentsettings.PaymentSettings, error)
}
