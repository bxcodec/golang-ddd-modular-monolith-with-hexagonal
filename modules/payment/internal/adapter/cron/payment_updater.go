package cron

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/modules/payment"
)

// PaymentUpdaterConfig contains configuration for the payment updater cron job
type PaymentUpdaterConfig struct {
	BatchSize int
	DryRun    bool
}

// PaymentUpdater is the cron adapter for payment update operations
type PaymentUpdater struct {
	paymentService payment.IPaymentService
	config         PaymentUpdaterConfig
}

// NewPaymentUpdater creates a new payment updater cron adapter
func NewPaymentUpdater(paymentService payment.IPaymentService, config PaymentUpdaterConfig) *PaymentUpdater {
	return &PaymentUpdater{
		paymentService: paymentService,
		config:         config,
	}
}

// ExecutionResult contains the results of a cron job execution
type ExecutionResult struct {
	StartTime      time.Time
	EndTime        time.Time
	Duration       time.Duration
	ProcessedCount int
	SuccessCount   int
	ErrorCount     int
	Errors         []error
}

// Execute runs the payment update cron job
func (u *PaymentUpdater) Execute() (resultData interface{}, err error) {
	result := &ExecutionResult{
		StartTime: time.Now(),
		Errors:    make([]error, 0),
	}

	log.Info().
		Str("started_at", result.StartTime.Format(time.RFC3339)).
		Msg("Payment update cron job started")

	if u.config.DryRun {
		log.Info().Msg("Running in DRY-RUN mode. No actual updates will be performed")
	}

	log.Info().
		Int("batch_size", u.config.BatchSize).
		Msg("Fetching pending payments")

	payments, _, err := u.paymentService.FetchPayments(payment.FetchPaymentsParams{
		Limit:  u.config.BatchSize,
		Status: "pending",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}

	log.Info().
		Int("count", len(payments)).
		Msg("Found payments to process")

	if len(payments) == 0 {
		log.Info().Msg("No payments to process")
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Process payments
	for i, p := range payments {
		if i >= u.config.BatchSize {
			log.Info().
				Int("batch_size", u.config.BatchSize).
				Msg("Batch size limit reached. Stopping processing")
			break
		}

		result.ProcessedCount++

		log.Info().
			Int("current", i+1).
			Int("total", len(payments)).
			Str("payment_id", p.ID).
			Str("status", p.Status).
			Float64("amount", p.Amount).
			Str("currency", p.Currency).
			Msg("Processing payment")

		// Apply business logic for payment updates
		if err := u.processPayment(&p); err != nil {
			log.Error().
				Err(err).
				Str("payment_id", p.ID).
				Msg("Failed to process payment")
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Errorf("payment %s: %w", p.ID, err))
			continue
		}

		result.SuccessCount++
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Print summary
	u.printSummary(result)

	if result.ErrorCount > 0 {
		return result, fmt.Errorf("cron job completed with %d error(s)", result.ErrorCount)
	}

	return result, nil
}

// processPayment handles the business logic for a single payment
func (u *PaymentUpdater) processPayment(p *payment.Payment) (err error) {
	// Skip payments that don't need processing
	if p.Status != "pending" {
		log.Info().
			Str("payment_id", p.ID).
			Str("status", p.Status).
			Msg("Skipped payment. Status is not pending")
		return nil
	}

	// Handle dry-run mode
	if u.config.DryRun {
		log.Info().
			Str("payment_id", p.ID).
			Msg("DRY-RUN mode. Would update payment to processing")
		return nil
	}

	// Update payment status
	p.Status = "processing"
	if err = u.paymentService.UpdatePayment(p); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	log.Info().
		Str("payment_id", p.ID).
		Str("new_status", "processing").
		Msg("Updated payment status")
	return nil
}

// printSummary prints the execution summary
func (u *PaymentUpdater) printSummary(result *ExecutionResult) {
	if result.ErrorCount > 0 {
		log.Warn().
			Str("completed_at", result.EndTime.Format(time.RFC3339)).
			Dur("duration", result.Duration).
			Int("processed", result.ProcessedCount).
			Int("success", result.SuccessCount).
			Int("errors", result.ErrorCount).
			Str("status", u.getJobStatus(result.ErrorCount)).
			Msg("Payment update cron job completed")
		return
	}

	log.Info().
		Str("completed_at", result.EndTime.Format(time.RFC3339)).
		Dur("duration", result.Duration).
		Int("processed", result.ProcessedCount).
		Int("success", result.SuccessCount).
		Int("errors", result.ErrorCount).
		Str("status", u.getJobStatus(result.ErrorCount)).
		Msg("Payment update cron job completed")
}

// getJobStatus returns a status string based on error count
func (u *PaymentUpdater) getJobStatus(errorCount int) (status string) {
	if errorCount == 0 {
		return "SUCCESS"
	}
	return "COMPLETED WITH ERRORS"
}
