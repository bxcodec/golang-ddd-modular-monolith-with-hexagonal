package paymentsettings

import "time"

// PaymentSettings represents payment configuration settings - your domain model
type PaymentSetting struct {
	ID           string    `json:"id"`
	SettingKey   string    `json:"settingKey"`
	SettingValue string    `json:"settingValue"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PaymentSettingFetchParams struct {
	Currency   string `json:"currency"`
	SettingKey string `json:"settingKey"`
	Limit      int    `json:"limit"`
	Cursor     string `json:"cursor"`
	Status     string `json:"status"`
}

// IPaymentSettingsService defines the public API for payment settings operations
type IPaymentSettingsService interface {
	FetchPaymentSettings(params PaymentSettingFetchParams) (result []PaymentSetting, nextCursor string, err error)
	CreatePaymentSetting(settings *PaymentSetting) error
	GetPaymentSetting(id string) (PaymentSetting, error)
	UpdatePaymentSetting(settings *PaymentSetting) error
	DeletePaymentSetting(id string) error
}
