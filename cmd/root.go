package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/config"
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	db      *sql.DB
)

var rootCmd = &cobra.Command{
	Use:   "payment-app",
	Short: "Payment application with DDD and Hexagonal Architecture",
	Long: `A modular monolith payment application built with Domain-Driven Design (DDD)
and Hexagonal Architecture principles.

This application supports multiple execution modes:
  - REST API server for handling HTTP requests
  - Cron jobs for scheduled payment updates`,
	PersistentPreRunE:  initApp,
	PersistentPostRunE: cleanupApp,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Application execution failed")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".env", "Config file path (supports .env files)")
}

func initApp(cmd *cobra.Command, args []string) (err error) {
	cfg, err = config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logger.Init(logger.Config{
		Level:       cfg.App.LogLevel,
		Environment: cfg.App.Environment,
	})

	log.Info().
		Str("app", cfg.App.Name).
		Str("environment", cfg.App.Environment).
		Str("log_level", cfg.App.LogLevel).
		Msg("Starting application")

	db, err = sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("host", cfg.Database.Host).
		Str("port", cfg.Database.Port).
		Str("database", cfg.Database.Name).
		Int("max_open_conns", cfg.Database.MaxOpenConns).
		Int("max_idle_conns", cfg.Database.MaxIdleConns).
		Msg("Successfully connected to PostgreSQL database")

	return nil
}

func cleanupApp(cmd *cobra.Command, args []string) (err error) {
	if db != nil {
		log.Info().Msg("Closing database connection")
		return db.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return db
}

func GetConfig() *config.Config {
	return cfg
}
