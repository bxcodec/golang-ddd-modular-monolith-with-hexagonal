package ports

import "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"

//go:generate mockery --name IPaymentRepository

type IPaymentRepository interface {
	CreatePayment(p *payment.Payment) error
	GetPayment(id string) (*payment.Payment, error)
	GetPayments() ([]*payment.Payment, error)
	UpdatePayment(p *payment.Payment) error
	DeletePayment(id string) error
}
