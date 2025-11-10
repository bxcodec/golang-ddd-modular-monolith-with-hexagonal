// Package ports defines the interfaces (ports) for external dependencies in hexagonal architecture.
//
// Ports are contracts that adapters must implement. By defining ports in the module:
//   - The core domain remains independent of external systems (databases, APIs, etc.)
//   - Adapters (implementations) can be swapped without changing business logic
//   - Testing becomes easier through mock implementations
//
// This follows the dependency inversion principle: high-level modules (domain) don't depend
// on low-level modules (adapters), both depend on abstractions (ports).
package ports

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

// IPaymentSettingsRepository is an outbound port for payment settings data persistence.
// This interface is defined by the domain and implemented by the repository adapter.
// The domain doesn't know or care about the underlying storage mechanism.
type IPaymentSettingsRepository interface {
	FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (result []paymentsettings.PaymentSetting, nextCursor string, err error)
	GetPaymentSetting(id string) (paymentsettings.PaymentSetting, error)
	CreatePaymentSetting(settings *paymentsettings.PaymentSetting) error
	UpdatePaymentSetting(settings *paymentsettings.PaymentSetting) error
	DeletePaymentSetting(id string) error
}
