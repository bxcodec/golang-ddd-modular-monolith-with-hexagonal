package cmd

import (
	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/spf13/cobra"
)

var (
	// Cron job configuration
	batchSize int
	dryRun    bool
)

// cronUpdatePaymentCmd represents the cron-update-payment command
var cronUpdatePaymentCmd = &cobra.Command{
	Use:   "cron-update-payment",
	Short: "Run payment update cron job",
	Long: `Run the payment update cron job for scheduled payment processing.

This command is designed to be executed by a cron scheduler to perform
batch updates on payments. The actual business logic resides in the
payment module's cron adapter (hexagonal architecture).

The command will:
  1. Initialize the payment module with the cron adapter
  2. Execute the cron adapter which handles all business logic
  3. Return the execution results

Example:
  # Run in production
  engine cron-update-payment

  # Dry run (no actual updates)
  engine cron-update-payment --dry-run

  # Process specific batch size
  engine cron-update-payment --batch-size 100

Cron schedule example (runs every hour):
  0 * * * * /path/to/engine cron-update-payment >> /var/log/payment-cron.log 2>&1`,
	RunE: runCronUpdatePayment,
}

func init() {
	rootCmd.AddCommand(cronUpdatePaymentCmd)

	// Cron-specific flags
	cronUpdatePaymentCmd.Flags().IntVar(&batchSize, "batch-size", 50, "Number of payments to process in one batch")
	cronUpdatePaymentCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Run in dry-run mode (no actual updates)")
}

func runCronUpdatePayment(cmd *cobra.Command, args []string) (err error) {
	// Get database connection from root command
	db := GetDB()

	// Initialize modules with cron configuration
	paymentSettingsModule := settingsfactory.NewModule(settingsfactory.ModuleConfig{
		DB: db,
	})

	paymentModule := paymentfactory.NewModule(paymentfactory.ModuleConfig{
		DB:                  db,
		PaymentSettingsPort: paymentSettingsModule.Service,
		CronBatchSize:       batchSize,
		CronDryRun:          dryRun,
	})

	// Execute the cron job via the adapter
	// The adapter contains all the business logic
	_, err = paymentModule.PaymentUpdater.Execute()

	return err
}
