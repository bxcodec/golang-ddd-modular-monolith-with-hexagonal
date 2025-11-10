package service

import (
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	paymentsettings "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports"
)

type PaymentService struct {
	paymentRepo         ports.IPaymentRepository
	paymentSettingsRepo ports.IPaymentSettingsPort
}

func NewPaymentService(paymentRepo ports.IPaymentRepository, paymentSettingsRepo ports.IPaymentSettingsPort) (service *PaymentService) {
	return &PaymentService{
		paymentRepo:         paymentRepo,
		paymentSettingsRepo: paymentSettingsRepo,
	}
}

func (s *PaymentService) CreatePayment(p *payment.Payment) (err error) {
	paymentSettings, _, err := s.paymentSettingsRepo.FetchPaymentSettings(paymentsettings.PaymentSettingFetchParams{
		Currency: p.Currency,
		Limit:    1,
		Cursor:   "",
	})
	if err != nil {
		return err
	}
	_ = paymentSettings
	return s.paymentRepo.CreatePayment(p)
}

func (s *PaymentService) GetPayment(id string) (result payment.Payment, err error) {
	p, err := s.paymentRepo.GetPayment(id)
	if err != nil {
		return payment.Payment{}, err
	}
	return p, nil
}

func (s *PaymentService) FetchPayments(params payment.FetchPaymentsParams) (result []payment.Payment, nextCursor string, err error) {
	return s.paymentRepo.FetchPayments(params)
}

func (s *PaymentService) UpdatePayment(p *payment.Payment) (err error) {
	return s.paymentRepo.UpdatePayment(p)
}

func (s *PaymentService) DeletePayment(id string) (err error) {
	return s.paymentRepo.DeletePayment(id)
}
