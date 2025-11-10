package paymentsettings

import "time"

// PaymentSettings represents payment configuration settings - your domain model
type PaymentSetting struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentSettingFetchParams struct {
	Currency string
	Limit    int
	Cursor   string
}

//go:generate mockery --name IPaymentSettingsService

// IPaymentSettingsService defines the public API for payment settings operations
type IPaymentSettingsService interface {
	FetchPaymentSettings(params PaymentSettingFetchParams) (res []PaymentSetting, nextCursor string, err error)
	CreatePaymentSetting(settings *PaymentSetting) error
	UpdatePaymentSetting(settings *PaymentSetting) error
	DeletePaymentSetting(id string) error
}
