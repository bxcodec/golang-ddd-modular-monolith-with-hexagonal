package dto

import (
	"time"

	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
)

type PaymentSettingResponse struct {
	ID           string    `json:"id"`
	SettingKey   string    `json:"settingKey"`
	SettingValue string    `json:"settingValue"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func FromPaymentSettingToResponse(setting paymentsettings.PaymentSetting) PaymentSettingResponse {
	return PaymentSettingResponse{
		ID:           setting.ID,
		SettingKey:   setting.SettingKey,
		SettingValue: setting.SettingValue,
		Currency:     setting.Currency,
		Status:       setting.Status,
		CreatedAt:    setting.CreatedAt,
		UpdatedAt:    setting.UpdatedAt,
	}
}

func FromPaymentSettingListToResponse(settings []paymentsettings.PaymentSetting) []PaymentSettingResponse {
	response := make([]PaymentSettingResponse, len(settings))
	for i, setting := range settings {
		response[i] = FromPaymentSettingToResponse(setting)
	}
	return response
}
