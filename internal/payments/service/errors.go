package service

import (
	"fmt"

	"github.com/gofrs/uuid"
)

type CreatePaymentPlanError struct{}

func (cp CreatePaymentPlanError) Error() string {
	return "failed to create payment plan"
}

type ListPaymentPlansByUserIDError struct {
	userID uuid.UUID
}

func (lp ListPaymentPlansByUserIDError) Error() string {
	return fmt.Sprintf("failed to get payment plans for user: %v", lp.userID)
}

type ListPaymentInstallmentsByPlanIDError struct {
	planID uuid.UUID
}

type CreatePaymentInstallmentError struct{}

func (cp CreatePaymentInstallmentError) Error() string {
	return "failed to create payment plan installment"
}

func (lp ListPaymentInstallmentsByPlanIDError) Error() string {
	return fmt.Sprintf("failed to get payment plans for user: %v", lp.planID)
}

type PaymentRecordNotFoundError struct {
	planID uuid.UUID
}

func (pr PaymentRecordNotFoundError) Error() string {
	return fmt.Sprintf("failed to get payment plan: %v", pr.planID)
}
