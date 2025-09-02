package payment

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"mobigo-backend/internal/installment"
	"time"
)

type Service interface {
	GenerateInstallmentPlan(ctx context.Context, req GeneratePlanRequest) error
}

type service struct {
	paymentRepo     Repository
	installmentRepo installment.Repository
	contextTimeout  time.Duration
}

func NewService(paymentRepo Repository, installmentRepo installment.Repository, timeout time.Duration) Service {
	return &service{
		paymentRepo:     paymentRepo,
		installmentRepo: installmentRepo,
		contextTimeout:  timeout,
	}
}

// GeneratePlanRequest holds the data needed to create a payment plan.
type GeneratePlanRequest struct {
	AgreementID        int64
	TotalPrice         float64
	DownPayment        float64
	Tenor              int     // in months
	AnnualInterestRate float64 // as a percentage, e.g., 12 for 12%
}

func (s *service) GenerateInstallmentPlan(ctx context.Context, req GeneratePlanRequest) error {
	// --- Validation ---
	if req.AgreementID == 0 || req.TotalPrice <= 0 || req.DownPayment < 0 || req.Tenor <= 0 {
		return errors.New("invalid request parameters")
	}
	if req.DownPayment >= req.TotalPrice {
		return errors.New("down payment must be less than the total price")
	}

	// --- Down Payment Record ---
	dpPayment := &domain.Payment{
		AgreementID:   req.AgreementID,
		Amount:        req.DownPayment,
		PaymentMethod: "Down Payment",
		Status:        domain.PaymentStatusPending,
	}
	if err := s.paymentRepo.CreatePayment(ctx, dpPayment); err != nil {
		return err
	}

	// --- Installment Calculation (Flat Rate) ---
	loanPrincipal := req.TotalPrice - req.DownPayment
	totalInterest := loanPrincipal * (req.AnnualInterestRate / 100) * (float64(req.Tenor) / 12)
	totalRepayment := loanPrincipal + totalInterest
	monthlyBill := totalRepayment / float64(req.Tenor)

	// --- Installment Plan Container Record ---
	installmentPayment := &domain.Payment{
		AgreementID:   req.AgreementID,
		Amount:        totalRepayment,
		PaymentMethod: "Installment",
		Status:        domain.PaymentStatusPending,
	}
	if err := s.paymentRepo.CreatePayment(ctx, installmentPayment); err != nil {
		return err
	}

	// --- Generate Individual Installment Records ---
	var installmentsToCreate []*domain.Installment
	for i := 1; i <= req.Tenor; i++ {
		dueDate := time.Now().AddDate(0, i, 0) // Simple due date logic: +i months from today
		inst := &domain.Installment{
			PaymentID:     installmentPayment.ID,
			DueDate:       dueDate,
			AmountDue:     monthlyBill,
			PenaltyAmount: 0,
			TotalDue:      monthlyBill,
			Status:        domain.InstallmentStatusPending,
		}
		installmentsToCreate = append(installmentsToCreate, inst)
	}

	return s.installmentRepo.CreateInstallments(ctx, installmentsToCreate)
}
