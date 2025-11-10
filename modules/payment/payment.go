package payment

import "time"

// Payment represents a payment transaction - your domain model
type Payment struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FetchPaymentsParams struct {
	Cursor   string `json:"cursor"`
	Limit    int    `json:"limit"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

// IPaymentService defines the public API for payment operations
type IPaymentService interface {
	CreatePayment(payment *Payment) (err error)
	GetPayment(id string) (payment Payment, err error)
	FetchPayments(params FetchPaymentsParams) (result []Payment, nextCursor string, err error)
	UpdatePayment(payment *Payment) (err error)
	DeletePayment(id string) (err error)
}
