package repo

import (
	"context"

	"golangreferenceapi/internal/payments"
)

//go:generate mockgen -source=./repository.go -destination=./mockrepository.go -package=repo
type Repository interface {
	CreatePaymentPlan(ctx context.Context, arg *payments.CreatePlanParams) (*payments.Plan, error)
	ListPaymentPlansByUserID(ctx context.Context, userID string) ([]*payments.Plan, error)
}
