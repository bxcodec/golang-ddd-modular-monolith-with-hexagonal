package paymentsettings

import "github.com/labstack/echo/v4"

// Module represents the payment-settings module with all its components
type Module struct {
	Service            IPaymentSettingsService
	RegisterController func(*echo.Group)
}

func (m *Module) RegisterHTTPHandlers(e *echo.Group) {
	if m.RegisterController != nil {
		m.RegisterController(e)
	}
}
