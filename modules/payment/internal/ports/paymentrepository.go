package ports

import "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"

type IPaymentRepository interface {
	CreatePayment(p *payment.Payment) error
	GetPayment(id string) (payment.Payment, error)
	FetchPayments(params payment.FetchPaymentsParams) (payments []payment.Payment, nextCursor string, err error)
	UpdatePayment(p *payment.Payment) error
	DeletePayment(id string) error
}
