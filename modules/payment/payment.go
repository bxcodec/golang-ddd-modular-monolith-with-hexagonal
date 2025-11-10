package payment

import "time"

// Payment represents a payment transaction - your domain model
type Payment struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//go:generate mockery --name IPaymentService

// IPaymentService defines the public API for payment operations
type IPaymentService interface {
	CreatePayment(payment *Payment) (err error)
	GetPayment(id string) (payment Payment, err error)
	GetPayments() (result []Payment, nextCursor string, err error)
	UpdatePayment(payment *Payment) (err error)
	DeletePayment(id string) (err error)
}
