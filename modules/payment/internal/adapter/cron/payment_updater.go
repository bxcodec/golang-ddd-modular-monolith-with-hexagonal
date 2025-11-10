package cron

import (
	"fmt"
	"log"
	"time"

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

	log.Printf("=== Payment Update Cron Job Started at %s ===\n", result.StartTime.Format(time.RFC3339))

	if u.config.DryRun {
		log.Println("Running in DRY-RUN mode - no actual updates will be performed")
	}

	// Fetch pending payments
	log.Printf("Fetching payments (batch size: %d)...\n", u.config.BatchSize)
	payments, _, err := u.paymentService.FetchPayments(payment.FetchPaymentsParams{
		Limit:  u.config.BatchSize,
		Status: "pending",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}

	log.Printf("Found %d payment(s) to process\n", len(payments))

	if len(payments) == 0 {
		log.Println("No payments to process")
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		return result, nil
	}

	// Process payments
	for i, p := range payments {
		if i >= u.config.BatchSize {
			log.Printf("Batch size limit reached (%d), stopping processing\n", u.config.BatchSize)
			break
		}

		result.ProcessedCount++

		log.Printf("[%d/%d] Processing payment ID: %s (Status: %s, Amount: %.2f %s)\n",
			i+1, len(payments), p.ID, p.Status, p.Amount, p.Currency)

		// Apply business logic for payment updates
		if err := u.processPayment(&p); err != nil {
			log.Printf("ERROR: Failed to process payment %s: %v\n", p.ID, err)
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
		log.Printf("INFO: Skipped payment %s (status: %s)\n", p.ID, p.Status)
		return nil
	}

	// Handle dry-run mode
	if u.config.DryRun {
		log.Printf("INFO: [DRY-RUN] Would update payment %s to processing\n", p.ID)
		return nil
	}

	// Update payment status
	p.Status = "processing"
	if err = u.paymentService.UpdatePayment(p); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	log.Printf("INFO: Updated payment %s to processing\n", p.ID)
	return nil
}

// printSummary prints the execution summary
func (u *PaymentUpdater) printSummary(result *ExecutionResult) {
	log.Println("\n=== Payment Update Cron Job Summary ===")
	log.Printf("Duration: %s\n", result.Duration)
	log.Printf("Processed: %d payment(s)\n", result.ProcessedCount)
	log.Printf("Success: %d\n", result.SuccessCount)
	if result.ErrorCount > 0 {
		log.Printf("Errors: %d\n", result.ErrorCount)
	}
	log.Printf("Status: %s\n", u.getJobStatus(result.ErrorCount))
	log.Printf("=== Cron Job Completed at %s ===\n", result.EndTime.Format(time.RFC3339))
}

// getJobStatus returns a status string based on error count
func (u *PaymentUpdater) getJobStatus(errorCount int) (status string) {
	if errorCount == 0 {
		return "SUCCESS"
	}
	return "COMPLETED WITH ERRORS"
}
