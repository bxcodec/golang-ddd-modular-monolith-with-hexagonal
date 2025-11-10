package dto

import (
	"time"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
)

type PaymentResponse struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
