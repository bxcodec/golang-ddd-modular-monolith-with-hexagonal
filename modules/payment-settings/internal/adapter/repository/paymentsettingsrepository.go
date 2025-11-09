package repository

import (
	"database/sql"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type PaymentSettingsRepository struct {
	db *sql.DB
}

func NewPaymentSettingsRepository(db *sql.DB) *PaymentSettingsRepository {
	return &PaymentSettingsRepository{db: db}
}

func (r *PaymentSettingsRepository) FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (res []paymentsettings.PaymentSetting, nextCursor string, err error) {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) GetPaymentSettingByCurrency(currency string) (paymentsettings.PaymentSetting, error) {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) CreatePaymentSetting(settings *paymentsettings.PaymentSetting) error {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) UpdatePaymentSetting(settings *paymentsettings.PaymentSetting) error {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) DeletePaymentSetting(id string) error {
	panic("not implemented")
}
