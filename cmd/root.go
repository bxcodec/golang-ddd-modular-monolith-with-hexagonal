package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	// Database configuration flags
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	dbSSLMode  string

	// Shared database connection
	db *sql.DB
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "payment-app",
	Short: "Payment application with DDD and Hexagonal Architecture",
	Long: `A modular monolith payment application built with Domain-Driven Design (DDD)
and Hexagonal Architecture principles.

This application supports multiple execution modes:
  - REST API server for handling HTTP requests
  - Cron jobs for scheduled payment updates`,
	PersistentPreRunE:  initDatabase,
	PersistentPostRunE: closeDatabase,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Database configuration flags
	rootCmd.PersistentFlags().StringVar(&dbHost, "db-host", getEnv("POSTGRES_HOST", "127.0.0.1"), "Database host")
	rootCmd.PersistentFlags().StringVar(&dbPort, "db-port", getEnv("POSTGRES_PORT", "5432"), "Database port")
	rootCmd.PersistentFlags().StringVar(&dbUser, "db-user", getEnv("POSTGRES_USER", "user"), "Database user")
	rootCmd.PersistentFlags().StringVar(&dbPassword, "db-password", getEnv("POSTGRES_PASSWORD", "password"), "Database password")
	rootCmd.PersistentFlags().StringVar(&dbName, "db-name", getEnv("POSTGRES_DB", "payment"), "Database name")
	rootCmd.PersistentFlags().StringVar(&dbSSLMode, "db-sslmode", getEnv("POSTGRES_SSLMODE", "disable"), "Database SSL mode")
}

// initDatabase initializes the database connection
func initDatabase(cmd *cobra.Command, args []string) (err error) {
	// Create PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Initialize database
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return nil
}

// closeDatabase closes the database connection
func closeDatabase(cmd *cobra.Command, args []string) (err error) {
	if db != nil {
		log.Println("Closing database connection...")
		return db.Close()
	}
	return nil
}

// GetDB returns the shared database connection
func GetDB() (database *sql.DB) {
	return db
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key string, defaultValue string) (value string) {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
