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

func NewPaymentService(paymentRepo ports.IPaymentRepository,
	paymentSettingsRepo ports.IPaymentSettingsPort) *PaymentService {
	return &PaymentService{
		paymentRepo:         paymentRepo,
		paymentSettingsRepo: paymentSettingsRepo,
	}
}

func (s *PaymentService) CreatePayment(p *payment.Payment) error {
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

func (s *PaymentService) GetPayment(id string) (payment.Payment, error) {
	p, err := s.paymentRepo.GetPayment(id)
	if err != nil {
		return payment.Payment{}, err
	}
	return *p, nil
}

func (s *PaymentService) GetPayments() (result []payment.Payment, nextCursor string, err error) {
	payments, err := s.paymentRepo.GetPayments()
	if err != nil {
		return nil, "", err
	}
	result = make([]payment.Payment, len(payments))
	for i, p := range payments {
		result[i] = *p
	}
	return result, "", nil
}

func (s *PaymentService) UpdatePayment(p *payment.Payment) error {
	return s.paymentRepo.UpdatePayment(p)
}

func (s *PaymentService) DeletePayment(id string) error {
	return s.paymentRepo.DeletePayment(id)
}
