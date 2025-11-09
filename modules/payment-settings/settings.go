package paymentsettings

import "time"

// PaymentSettings represents payment configuration settings - your domain model
type PaymentSettings struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IPaymentSettingsService defines the public API for payment settings operations
type IPaymentSettingsService interface {
	GetPaymentSettingsByCurrency(currency string) (PaymentSettings, error)
	CreatePaymentSettings(paymentSettings *PaymentSettings) error
	UpdatePaymentSettings(paymentSettings *PaymentSettings) error
	DeletePaymentSettings(id string) error
}
