package ports

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type IPaymentSettingsRepository interface {
	FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (result []paymentsettings.PaymentSetting, nextCursor string, err error)
	GetPaymentSetting(id string) (paymentsettings.PaymentSetting, error)
	CreatePaymentSetting(settings *paymentsettings.PaymentSetting) error
	UpdatePaymentSetting(settings *paymentsettings.PaymentSetting) error
	DeletePaymentSetting(id string) error
}
