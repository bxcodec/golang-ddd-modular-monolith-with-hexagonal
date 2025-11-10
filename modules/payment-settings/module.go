package paymentsettings

import "github.com/labstack/echo/v4"

// Module encapsulates the Payment Settings module following hexagonal architecture.
//
// Structure:
//   - Service: The hexagon core containing business logic
//   - RegisterController: Inbound adapter for HTTP/REST API
//
// This module is self-contained and can be composed with other modules in the monolith.
// All dependencies are injected via the factory, maintaining loose coupling and testability.
type Module struct {
	Service            IPaymentSettingsService
	RegisterController func(*echo.Group)
}

// RegisterHTTPHandlers registers all HTTP endpoints for this module.
// This allows the monolith to compose multiple modules by registering their routes
// into the main HTTP router, maintaining module encapsulation.
func (m *Module) RegisterHTTPHandlers(e *echo.Group) {
	if m.RegisterController != nil {
		m.RegisterController(e)
	}
}
