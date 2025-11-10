package ports

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

// IPaymentSettingsPort defines the contract for payment settings that the payment module needs.
// If you see theres a less method here compared to paymentsettings.IPaymentSettingsService, it's because we don't need all the methods.
// We only need the ones that the payment module needs.
// This just to follow Golang idiomatic way of doing things.
// Interface should be used to abstract external dependencies. And only define the necessary methods.
type IPaymentSettingsPort interface {
	FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (res []paymentsettings.PaymentSetting, nextCursor string, err error)
}
