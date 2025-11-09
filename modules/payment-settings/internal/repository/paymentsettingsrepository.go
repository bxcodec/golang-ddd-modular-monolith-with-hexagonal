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

func (r *PaymentSettingsRepository) GetPaymentSettingsByCurrency(currency string) (paymentsettings.PaymentSettings, error) {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) CreatePaymentSettings(settings *paymentsettings.PaymentSettings) error {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) UpdatePaymentSettings(settings *paymentsettings.PaymentSettings) error {
	panic("not implemented")
}

func (r *PaymentSettingsRepository) DeletePaymentSettings(id string) error {
	panic("not implemented")
}
