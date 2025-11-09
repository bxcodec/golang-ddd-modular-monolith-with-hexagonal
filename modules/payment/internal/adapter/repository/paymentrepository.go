package repository

import (
	"database/sql"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/internal/ports"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) ports.IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(p *payment.Payment) error {
	panic("not implemented")
}

func (r *PaymentRepository) GetPayment(id string) (*payment.Payment, error) {
	panic("not implemented")
}

func (r *PaymentRepository) GetPayments() ([]*payment.Payment, error) {
	panic("not implemented")
}

func (r *PaymentRepository) UpdatePayment(p *payment.Payment) error {
	panic("not implemented")
}

func (r *PaymentRepository) DeletePayment(id string) error {
	panic("not implemented")
}
