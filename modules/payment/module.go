package payment

import "github.com/labstack/echo/v4"

// CronAdapter defines the contract for scheduled job operations within the module.
// This is an inbound adapter allowing external cron schedulers to trigger module logic.
type CronAdapter interface {
	Execute() (interface{}, error)
}

// Module encapsulates the Payment module following hexagonal architecture.
//
// Structure:
//   - Service: The hexagon core (business logic)
//   - RegisterController: Inbound adapter (REST API)
//   - PaymentUpdater: Inbound adapter (cron job)
//
// The Module is the deployable unit in our modular monolith. It contains everything needed
// for payment operations: domain logic, HTTP handlers, scheduled jobs, and database access.
// All external dependencies are injected via ports to maintain module independence.
type Module struct {
	Service            IPaymentService
	RegisterController func(*echo.Group)
	// Cron adapters for scheduled jobs
	PaymentUpdater CronAdapter
}

// RegisterHTTPHandlers registers all HTTP endpoints for this module.
// This allows the monolith to compose multiple modules by registering their routes
// into the main HTTP router, maintaining module encapsulation.
func (m *Module) RegisterHTTPHandlers(e *echo.Group) {
	if m.RegisterController != nil {
		m.RegisterController(e)
	}
}
