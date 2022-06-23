package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRecordNotFound = errors.New("record not fund")
	ErrGenerateUUID   = errors.New("failed to generate uuid")
)

type PaymentPlanService interface {
	// GetPaymentPlanByUserID gets payment plans selected by userID
	GetPaymentPlanByUserID(ctx context.Context, userID uuid.UUID) ([]PaymentPlans, error)

	// CreatePendingPaymentPlan
	CreatePendingPaymentPlan(
		ctx context.Context,
		userID uuid.UUID,
		paymentPlan *CreatePaymentPlanParams,
	) (*PaymentPlans, error)

	// CompletePaymentPlanCreation
	CompletePaymentPlanCreation(
		ctx context.Context,
		userID uuid.UUID,
		paymentPlanID uuid.UUID,
	) error
}

type PaymentPlanInstallment struct {
	ID       string `json:"id"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	DueAt    string `json:"due_at"`
	Status   string `json:"status"`
}

type PaymentPlans struct {
	ID           string `json:"id"`
	Currency     string `json:"currency"`
	TotalAmount  string `json:"total_amount"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	Installments []PaymentPlanInstallment
}

type PaymentPlanInstallmentParams struct {
	ID       uuid.UUID `json:"id"`
	Amount   string    `json:"amount"`
	Currency string    `json:"currency"`
	DueAt    time.Time `json:"due_at"`
}

type CreatePaymentPlanParams struct {
	ID           uuid.UUID `json:"id"`
	Currency     string    `json:"currency"`
	TotalAmount  string    `json:"total_amount"`
	Installments []PaymentPlanInstallmentParams
}
