package ports

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

// IPaymentSettingsPort is an outbound port for inter-module communication with the Payment Settings module.
//
// Key architectural points:
//   - This interface has fewer methods than paymentsettings.IPaymentSettingsService
//   - We only define what the Payment module actually needs, not the entire API
//   - This prevents tight coupling and follows interface segregation principle
//   - The Payment module depends on this port, not directly on the Payment Settings module
//
// This demonstrates how modules communicate in a modular monolith: through well-defined ports,
// not direct dependencies. Each module defines its own view of what it needs from others.
//
// IMPORTANT - Circular Dependency Prevention:
// This port imports structs from the payment-settings module (PaymentSetting, PaymentSettingFetchParams).
// To prevent circular dependencies, establish a clear dependency direction:
//   - Payment module CAN import from Payment Settings module (current direction)
//   - Payment Settings module MUST NOT import from Payment module
//
// If bidirectional communication is needed in the future, consider:
//   - Creating a shared types package for common structs
//   - Using DTOs (Data Transfer Objects) instead of direct struct dependencies
type IPaymentSettingsPort interface {
	FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (res []paymentsettings.PaymentSetting, nextCursor string, err error)
}
