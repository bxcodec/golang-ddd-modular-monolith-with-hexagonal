package dto

import (
	"time"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type PaymentSettingResponse struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromPaymentSettingToResponse(setting paymentsettings.PaymentSetting) PaymentSettingResponse {
	return PaymentSettingResponse{
		ID:        setting.ID,
		Amount:    setting.Amount,
		Currency:  setting.Currency,
		Status:    setting.Status,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}

func FromPaymentSettingListToResponse(settings []paymentsettings.PaymentSetting) []PaymentSettingResponse {
	response := make([]PaymentSettingResponse, len(settings))
	for i, setting := range settings {
		response[i] = FromPaymentSettingToResponse(setting)
	}
	return response
}
