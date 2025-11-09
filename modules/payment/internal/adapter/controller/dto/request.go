package dto

import (
	"time"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
)

type CreatePaymentRequest struct {
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *CreatePaymentRequest) ToPayment() payment.Payment {
	return payment.Payment{
		Amount:    r.Amount,
		Currency:  r.Currency,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func FromPaymentToCreateRequest(p payment.Payment) CreatePaymentRequest {
	return CreatePaymentRequest{
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type UpdatePaymentRequest struct {
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *UpdatePaymentRequest) ToPayment(id string) payment.Payment {
	return payment.Payment{
		ID:        id,
		Amount:    r.Amount,
		Currency:  r.Currency,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
