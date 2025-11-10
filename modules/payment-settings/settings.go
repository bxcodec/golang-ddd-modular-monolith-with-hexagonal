// Package paymentsettings implements the Payment Settings module in a modular monolith architecture.
//
// This module manages payment configuration and settings following hexagonal architecture:
//   - Domain models (PaymentSetting) represent configuration entities
//   - Service interface (IPaymentSettingsService) exposes the module's public API
//   - Other modules consume this API through ports to access payment settings
//
// Module independence is maintained through interface-based communication,
// allowing this module to evolve without impacting dependent modules.
package paymentsettings

import "time"

// PaymentSetting represents payment configuration in the domain model.
// This is the core entity for managing payment-related settings and configurations.
type PaymentSetting struct {
	ID           string    `json:"id"`
	SettingKey   string    `json:"settingKey"`
	SettingValue string    `json:"settingValue"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// PaymentSettingFetchParams contains filtering and pagination parameters for querying payment settings.
type PaymentSettingFetchParams struct {
	Currency   string `json:"currency"`
	SettingKey string `json:"settingKey"`
	Limit      int    `json:"limit"`
	Cursor     string `json:"cursor"`
	Status     string `json:"status"`
}

// IPaymentSettingsService defines the public API of the Payment Settings module.
// This interface represents the module's full capabilities, but other modules should NOT import this directly.
// Instead, other modules define their own port interfaces (like IPaymentSettingsPort in the payment module)
// specifying only the methods they need. This maintains loose coupling and follows interface segregation.
type IPaymentSettingsService interface {
	FetchPaymentSettings(params PaymentSettingFetchParams) (result []PaymentSetting, nextCursor string, err error)
	CreatePaymentSetting(settings *PaymentSetting) error
	GetPaymentSetting(id string) (PaymentSetting, error)
	UpdatePaymentSetting(settings *PaymentSetting) error
	DeletePaymentSetting(id string) error
}
