// Package payment implements the Payment module in a modular monolith architecture.
//
// This module follows hexagonal architecture principles where:
//   - Domain models (Payment) represent the core business logic
//   - Service interfaces (IPaymentService) define the module's public API
//   - Ports define contracts for external dependencies
//   - Adapters implement the ports (repositories, controllers, cron jobs)
//
// Module boundaries are enforced through well-defined interfaces, allowing this module
// to operate independently while communicating with other modules via ports.
package payment

import "time"

// Payment represents a payment transaction in the domain model.
// This is the core entity in the payment bounded context.
type Payment struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// FetchPaymentsParams contains filtering and pagination parameters for querying payments.
type FetchPaymentsParams struct {
	Cursor   string `json:"cursor"`
	Limit    int    `json:"limit"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

// IPaymentService defines the public API of the Payment module.
// This is the primary interface exposed to other modules in the monolith,
// forming the module boundary. Other modules depend on this interface, not the implementation.
type IPaymentService interface {
	CreatePayment(payment *Payment) (err error)
	GetPayment(id string) (payment Payment, err error)
	FetchPayments(params FetchPaymentsParams) (result []Payment, nextCursor string, err error)
	UpdatePayment(payment *Payment) (err error)
	DeletePayment(id string) (err error)
}
