package payment

import "github.com/labstack/echo/v4"

//go:generate mockery --name CronAdapter

// CronAdapter represents the interface for cron job operations
type CronAdapter interface {
	Execute() (interface{}, error)
}

// Module represents the payment module with all its components
// This encapsulates the hexagonal core and its adapters
type Module struct {
	Service            IPaymentService
	RegisterController func(*echo.Group)
	// Cron adapters for scheduled jobs
	PaymentUpdater CronAdapter
}

// RegisterHTTPHandlers registers all HTTP endpoints for this module
// This method allows the module to expose its REST API
func (m *Module) RegisterHTTPHandlers(e *echo.Group) {
	if m.RegisterController != nil {
		m.RegisterController(e)
	}
}
