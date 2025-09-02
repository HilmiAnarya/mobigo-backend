package task

import (
	"context"
	"log"
	"mobigo-backend/internal/domain"
	"mobigo-backend/internal/installment"
	"time"
)

// PenaltyChecker contains the logic for checking and applying penalties.
type PenaltyChecker struct {
	installmentRepo installment.Repository
}

// NewPenaltyChecker creates a new instance of the PenaltyChecker.
func NewPenaltyChecker(repo installment.Repository) *PenaltyChecker {
	return &PenaltyChecker{
		installmentRepo: repo,
	}
}

// Run is the function that will be executed by the cron job.
func (pc *PenaltyChecker) Run() {
	log.Println("CRON JOB: Starting check for overdue installments...")

	// Define a daily penalty amount. In a real app, this should come from a config file.
	const dailyPenalty = 10000.0 // Rp 10,000

	ctx := context.Background()

	// 1. Get all overdue installments.
	overdueInstallments, err := pc.installmentRepo.FindOverdueInstallments(ctx)
	if err != nil {
		log.Printf("CRON ERROR: Could not fetch overdue installments: %v", err)
		return
	}

	if len(overdueInstallments) == 0 {
		log.Println("CRON JOB: No overdue installments found. Finished.")
		return
	}

	log.Printf("CRON JOB: Found %d overdue installment(s). Applying penalties...", len(overdueInstallments))

	// 2. Loop through each one and apply the penalty.
	for _, inst := range overdueInstallments {
		// Mark as overdue if it's currently pending
		if inst.Status == domain.InstallmentStatusPending {
			inst.Status = domain.InstallmentStatusOverdue
		}

		// Calculate how many days late the payment is.
		daysLate := int(time.Since(inst.DueDate).Hours() / 24)
		if daysLate < 1 {
			daysLate = 1 // Minimum 1 day late
		}

		// Calculate the new penalty amount.
		newPenalty := float64(daysLate) * dailyPenalty

		// Update the installment record
		inst.PenaltyAmount = newPenalty
		inst.TotalDue = inst.AmountDue + newPenalty
		inst.UpdatedAt = time.Now()

		// 3. Update the record in the database.
		if err := pc.installmentRepo.UpdateInstallment(ctx, inst); err != nil {
			log.Printf("CRON ERROR: Failed to update installment ID %d: %v", inst.ID, err)
			// Continue to the next installment even if this one fails
			continue
		}
		log.Printf("CRON JOB: Applied penalty to installment ID %d. Days late: %d, New Total Due: %.2f", inst.ID, daysLate, inst.TotalDue)
	}

	log.Println("CRON JOB: Finished applying penalties.")
}
