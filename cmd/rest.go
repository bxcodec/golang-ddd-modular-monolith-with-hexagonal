package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start the REST API server",
	Long: `Start the REST API server to handle HTTP requests.

This command initializes the application modules and starts an HTTP server
using the Echo framework with graceful shutdown support.

Example:
  payment-app rest
  payment-app rest --config .env.production`,
	RunE: runREST,
}

func init() {
	rootCmd.AddCommand(restCmd)
}

func runREST(cmd *cobra.Command, args []string) (err error) {
	cfg := GetConfig()
	db := GetDB()

	log.Info().Msg("Initializing REST API server")

	paymentSettingsModule := settingsfactory.NewModule(settingsfactory.ModuleConfig{
		DB: db,
	})

	paymentModule := paymentfactory.NewModule(paymentfactory.ModuleConfig{
		DB:                  db,
		PaymentSettingsPort: paymentSettingsModule.Service,
	})

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = middlewares.ErrorHandler

	if !cfg.IsProduction() {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(middlewares.CORS())
	e.Use(middlewares.SetRequestContextWithTimeout(cfg.Server.ReadTimeout))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":      "ok",
			"mode":        "rest",
			"environment": cfg.App.Environment,
		})
	})

	api := e.Group("/api/v1")
	paymentModule.RegisterHTTPHandlers(api)
	paymentSettingsModule.RegisterHTTPHandlers(api)

	go func() {
		log.Info().
			Str("port", cfg.Server.Port).
			Str("health_check", fmt.Sprintf("http://localhost:%s/health", cfg.Server.Port)).
			Str("api_base", fmt.Sprintf("http://localhost:%s/api/v1", cfg.Server.Port)).
			Msg("REST API server started")

		if err := e.Start(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
			log.Error().Err(err).Msg("Server stopped")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
		return err
	}

	log.Info().Msg("Server shutdown complete")
	return nil
}
