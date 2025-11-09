package dto

import (
	"time"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
)

type PaymentResponse struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromPaymentToResponse(p payment.Payment) PaymentResponse {
	return PaymentResponse{
		ID:        p.ID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type PaymentListResponse []PaymentResponse

func FromPaymentListToResponse(p []payment.Payment) PaymentListResponse {
	payments := make(PaymentListResponse, len(p))
	for i, p := range p {
		payments[i] = FromPaymentToResponse(p)
	}
	return payments
}
