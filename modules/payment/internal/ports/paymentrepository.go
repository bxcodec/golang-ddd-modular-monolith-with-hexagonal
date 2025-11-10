// Package ports defines the interfaces (ports) for external dependencies in hexagonal architecture.
//
// Ports are contracts that adapters must implement. By defining ports in the module:
//   - The core domain remains independent of external systems (databases, APIs, etc.)
//   - Adapters (implementations) can be swapped without changing business logic
//   - Testing becomes easier through mock implementations
//
// This follows the dependency inversion principle: high-level modules (domain) don't depend
// on low-level modules (adapters), both depend on abstractions (ports).
package ports

import "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"

// IPaymentRepository is an outbound port for payment data persistence.
// This interface is defined by the domain and implemented by the repository adapter.
// The domain doesn't know or care if this uses PostgreSQL, MongoDB, or in-memory storage.
type IPaymentRepository interface {
	CreatePayment(p *payment.Payment) error
	GetPayment(id string) (payment.Payment, error)
	FetchPayments(params payment.FetchPaymentsParams) (payments []payment.Payment, nextCursor string, err error)
	UpdatePayment(p *payment.Payment) error
	DeletePayment(id string) error
}
