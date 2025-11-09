package dto

import (
	"time"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type CreatePaymentSettingRequest struct {
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *CreatePaymentSettingRequest) ToPaymentSetting() paymentsettings.PaymentSetting {
	return paymentsettings.PaymentSetting{
		Amount:    r.Amount,
		Currency:  r.Currency,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

type UpdatePaymentSettingRequest struct {
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *UpdatePaymentSettingRequest) ToPaymentSetting(id string) paymentsettings.PaymentSetting {
	return paymentsettings.PaymentSetting{
		ID:        id,
		Amount:    r.Amount,
		Currency:  r.Currency,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
