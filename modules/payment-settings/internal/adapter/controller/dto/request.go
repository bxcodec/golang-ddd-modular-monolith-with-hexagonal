package dto

import (
	"time"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type CreatePaymentSettingRequest struct {
	SettingKey   string    `json:"settingKey"`
	SettingValue string    `json:"settingValue"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (r *CreatePaymentSettingRequest) ToPaymentSetting() paymentsettings.PaymentSetting {
	return paymentsettings.PaymentSetting{
		SettingKey:   r.SettingKey,
		SettingValue: r.SettingValue,
		Currency:     r.Currency,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

type UpdatePaymentSettingRequest struct {
	SettingKey   string    `json:"settingKey"`
	SettingValue string    `json:"settingValue"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (r *UpdatePaymentSettingRequest) ToPaymentSetting(id string) paymentsettings.PaymentSetting {
	return paymentsettings.PaymentSetting{
		ID:           id,
		SettingKey:   r.SettingKey,
		SettingValue: r.SettingValue,
		Currency:     r.Currency,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
