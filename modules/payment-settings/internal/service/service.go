package service

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/ports"
)

type PaymentSettingsService struct {
	repo ports.IPaymentSettingsRepository
}

func NewPaymentSettingsService(repo ports.IPaymentSettingsRepository) (service *PaymentSettingsService) {
	return &PaymentSettingsService{repo: repo}
}

func (s *PaymentSettingsService) CreatePaymentSetting(settings *paymentsettings.PaymentSetting) (err error) {
	return s.repo.CreatePaymentSetting(settings)
}

func (s *PaymentSettingsService) UpdatePaymentSetting(settings *paymentsettings.PaymentSetting) (err error) {
	return s.repo.UpdatePaymentSetting(settings)
}

func (s *PaymentSettingsService) DeletePaymentSetting(id string) (err error) {
	return s.repo.DeletePaymentSetting(id)
}

func (s *PaymentSettingsService) FetchPaymentSettings(params paymentsettings.PaymentSettingFetchParams) (res []paymentsettings.PaymentSetting, nextCursor string, err error) {
	return s.repo.FetchPaymentSettings(params)
}
