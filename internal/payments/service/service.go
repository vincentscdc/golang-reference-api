package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

//go:generate mockgen -source=./service.go -destination=../mock/servicemock/service_mock.go -package=servicemock
type PaymentPlanService interface {
	// GetPaymentPlanByUserID gets payment plans selected by userID
	GetPaymentPlanByUserID(ctx context.Context, userID uuid.UUID) ([]PaymentPlans, error)

	// CreatePendingPaymentPlan
	CreatePendingPaymentPlan(
		ctx context.Context,
		paymentPlan *CreatePaymentPlanParams,
	) (*PaymentPlans, error)

	// CompletePaymentPlanCreation
	CompletePaymentPlanCreation(
		ctx context.Context,
		paymentPlanID uuid.UUID,
		paymentPlan *CompletePaymentPlanParams,
	) (*PaymentPlans, error)
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
	UserID       string `json:"user_id"`
	Currency     string `json:"currency"`
	TotalAmount  string `json:"total_amount"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	Installments []PaymentPlanInstallment
}

type PaymentPlanInstallmentParams struct {
	Amount   string    `json:"amount"`
	Currency string    `json:"currency"`
	DueAt    time.Time `json:"due_at"`
}

type CreatePaymentPlanParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Currency     string    `json:"currency"`
	TotalAmount  string    `json:"total_amount"`
	Installments []PaymentPlanInstallmentParams
}

type CompletePaymentPlanParams struct {
	UserID uuid.UUID `json:"user_id"`
}
