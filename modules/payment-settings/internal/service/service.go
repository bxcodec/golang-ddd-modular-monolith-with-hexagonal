package service

import (
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/internal/ports"
)

type PaymentSettingsService struct {
	repo ports.IPaymentSettingsRepository
}

func NewPaymentSettingsService(repo ports.IPaymentSettingsRepository) *PaymentSettingsService {
	return &PaymentSettingsService{repo: repo}
}

func (s *PaymentSettingsService) GetPaymentSettingsByCurrency(currency string) (paymentsettings.PaymentSettings, error) {
	return s.repo.GetPaymentSettingsByCurrency(currency)
}

func (s *PaymentSettingsService) CreatePaymentSettings(settings *paymentsettings.PaymentSettings) error {
	return s.repo.CreatePaymentSettings(settings)
}

func (s *PaymentSettingsService) UpdatePaymentSettings(settings *paymentsettings.PaymentSettings) error {
	return s.repo.UpdatePaymentSettings(settings)
}

func (s *PaymentSettingsService) DeletePaymentSettings(id string) error {
	return s.repo.DeletePaymentSettings(id)
}
