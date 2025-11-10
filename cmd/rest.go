package cmd

import (
	"log"
	"time"

	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var (
	// REST server configuration
	restPort string
)

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start the REST API server",
	Long: `Start the REST API server to handle HTTP requests.

This command initializes the application modules and starts an HTTP server
using the Echo framework.

Example:
  engine rest --port 8080
  engine rest --db-host localhost --db-port 5432`,
	RunE: runREST,
}

func init() {
	rootCmd.AddCommand(restCmd)

	// REST-specific flags
	restCmd.Flags().StringVarP(&restPort, "port", "p", "8080", "Port to run the REST API server")
}

func runREST(cmd *cobra.Command, args []string) (err error) {
	log.Println("Starting REST API server...")

	// Get database connection from root command
	db := GetDB()

	// Initialize modules in dependency order (like Nest.js)
	// Factory handles ALL the wiring: hexagon core + adapters

	// 1. Payment Settings module (no dependencies)
	paymentSettingsModule := settingsfactory.NewModule(settingsfactory.ModuleConfig{
		DB: db,
	})

	// 2. Payment module (depends on Payment Settings)
	paymentModule := paymentfactory.NewModule(paymentfactory.ModuleConfig{
		DB:                  db,
		PaymentSettingsPort: paymentSettingsModule.Service,
	})

	// Setup HTTP server and register module routes
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.CORS())
	e.Use(middlewares.SetRequestContextWithTimeout(10 * time.Second))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
			"mode":   "rest",
		})
	})

	// Register module HTTP handlers (inbound adapters)
	api := e.Group("/api/v1")
	paymentModule.RegisterHTTPHandlers(api)
	paymentSettingsModule.RegisterHTTPHandlers(api)

	// Start server
	log.Printf("REST API server listening on port %s\n", restPort)
	log.Printf("Health check: http://localhost:%s/health\n", restPort)
	log.Printf("API endpoints: http://localhost:%s/api/v1\n", restPort)

	if err := e.Start(":" + restPort); err != nil {
		return err
	}

	return nil
}
