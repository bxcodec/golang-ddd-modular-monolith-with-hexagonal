package cmd

import (
	settingsfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment-settings/factory"
	paymentfactory "github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment/factory"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	batchSize int
	dryRun    bool
)

var cronUpdatePaymentCmd = &cobra.Command{
	Use:   "cron-update-payment",
	Short: "Run payment update cron job",
	Long: `Run the payment update cron job for scheduled payment processing.

This command is designed to be executed by a cron scheduler to perform
batch updates on payments. The actual business logic resides in the
payment module's cron adapter (hexagonal architecture).

Example:
  payment-app cron-update-payment
  payment-app cron-update-payment --dry-run
  payment-app cron-update-payment --batch-size 100

Cron schedule example (runs every hour):
  0 * * * * /path/to/payment-app cron-update-payment >> /var/log/payment-cron.log 2>&1`,
	RunE: runCronUpdatePayment,
}

func init() {
	rootCmd.AddCommand(cronUpdatePaymentCmd)

	cronUpdatePaymentCmd.Flags().IntVar(&batchSize, "batch-size", 0, "Number of payments to process in one batch (overrides CRON_BATCH_SIZE)")
	cronUpdatePaymentCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Run in dry-run mode (overrides CRON_DRY_RUN)")
}

func runCronUpdatePayment(cmd *cobra.Command, args []string) (err error) {
	cfg := GetConfig()
	db := GetDB()

	if batchSize == 0 {
		batchSize = cfg.Cron.BatchSize
	}
	if !dryRun {
		dryRun = cfg.Cron.DryRun
	}

	log.Info().
		Int("batch_size", batchSize).
		Bool("dry_run", dryRun).
		Msg("Starting payment update cron job")

	paymentSettingsModule := settingsfactory.NewModule(settingsfactory.ModuleConfig{
		DB: db,
	})

	paymentModule := paymentfactory.NewModule(paymentfactory.ModuleConfig{
		DB:                  db,
		PaymentSettingsPort: paymentSettingsModule.Service,
		CronBatchSize:       batchSize,
		CronDryRun:          dryRun,
	})

	result, err := paymentModule.PaymentUpdater.Execute()
	if err != nil {
		log.Error().Err(err).Msg("Cron job failed")
		return err
	}

	log.Info().Interface("result", result).Msg("Cron job completed successfully")
	return nil
}
