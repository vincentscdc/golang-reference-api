package repo

import (
	"context"

	"golangreferenceapi/internal/payments"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./repository.go -destination=./mockrepository.go -package=repo
type Repository interface {
	CreatePaymentPlan(ctx context.Context, arg *payments.CreatePlanParams) (*payments.Plan, error)
	ListPaymentPlansByUserID(ctx context.Context, userID uuid.UUID) ([]*payments.Plan, error)
	CreatePaymentInstallment(ctx context.Context, arg *payments.CreateInstallmentParams) (*payments.Installment, error)
	ListPaymentInstallmentsByPlanID(ctx context.Context, planID uuid.UUID) ([]*payments.Installment, error)
}
