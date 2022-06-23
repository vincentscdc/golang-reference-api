//nolint:cyclop // ignore for testing
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"golangreferenceapi/internal/payments/common"

	"github.com/google/uuid"
)

func TestPaymentServiceImp_CreatePendingPaymentPlan(t *testing.T) {
	t.Parallel()

	var (
		userID   = uuid.New()
		ctx      = context.Background()
		dueAt, _ = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          uuid.New(),
		Currency:    "usdc",
		TotalAmount: "1000",
		Installments: []PaymentPlanInstallmentParams{
			{
				Currency: "usdc",
				Amount:   "1000",
				DueAt:    dueAt,
			},
		},
	}

	service := NewPaymentPlanService()

	_, err := service.CreatePendingPaymentPlan(ctx, userID, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	payments, err := service.GetPaymentPlanByUserID(ctx, userID)
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	switch {
	case payments[0].ID != paymentPlanParams.ID.String():
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.ID, payments[0].ID)
	case payments[0].TotalAmount != paymentPlanParams.TotalAmount:
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.TotalAmount, payments[0].TotalAmount)
	case payments[0].Currency != paymentPlanParams.Currency:
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.Currency, payments[0].Currency)
	case payments[0].Status != paymentPlanStatusPending:
		t.Errorf("expected: %v, actual: %v", paymentPlanStatusPending, payments[0].Status)
	}

	installment := payments[0].Installments[0]
	installmentParams := paymentPlanParams.Installments[0]

	switch {
	case installment.Currency != installmentParams.Currency:
		t.Errorf("expected: %v, actual: %v",
			installment.Currency,
			installmentParams.Currency,
		)
	case installment.Amount != installmentParams.Amount:
		t.Errorf("expected: %v, actual: %v",
			installmentParams.Amount,
			installment.Amount,
		)
	case installment.Status != PaymentInstallmentStatusPending:
		t.Errorf("expected: %v, actual: %v",
			PaymentInstallmentStatusPending,
			installment.Status,
		)
	}

	if dueAt, err := time.Parse(common.TimeFormat, installment.DueAt); err != nil || dueAt.IsZero() {
		t.Errorf("unexpected %v", installment.DueAt)
	}

	if id, err := uuid.Parse(installment.ID); err != nil || id == uuid.Nil {
		t.Errorf("unexpected %v", installment.ID)
	}
}

func TestPaymentServiceImp_GetPaymentPlanByUserID(t *testing.T) {
	t.Parallel()

	var (
		userID   = uuid.New()
		ctx      = context.Background()
		dueAt, _ = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          uuid.New(),
		Currency:    "usdc",
		TotalAmount: "1000",
		Installments: []PaymentPlanInstallmentParams{
			{
				Currency: "usdc",
				Amount:   "1000",
				DueAt:    dueAt,
			},
		},
	}

	service := NewPaymentPlanService()

	_, err := service.CreatePendingPaymentPlan(ctx, userID, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	payments, err := service.GetPaymentPlanByUserID(ctx, uuid.New())
	if !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("unexpected record existed: %v", payments)
	}

	payments, err = service.GetPaymentPlanByUserID(ctx, userID)
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	switch {
	case payments[0].ID != paymentPlanParams.ID.String():
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.ID, payments[0].ID)
	case payments[0].TotalAmount != paymentPlanParams.TotalAmount:
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.TotalAmount, payments[0].TotalAmount)
	case payments[0].Currency != paymentPlanParams.Currency:
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.Currency, payments[0].Currency)
	case payments[0].Status != paymentPlanStatusPending:
		t.Errorf("expected: %v, actual: %v", paymentPlanStatusPending, payments[0].Status)
	}

	installment := payments[0].Installments[0]
	installmentParams := paymentPlanParams.Installments[0]

	switch {
	case installment.Currency != installmentParams.Currency:
		t.Errorf("expected: %v, actual: %v",
			installment.Currency,
			installmentParams.Currency,
		)
	case installment.Amount != installmentParams.Amount:
		t.Errorf("expected: %v, actual: %v",
			installmentParams.Amount,
			installment.Amount,
		)
	case installment.Status != PaymentInstallmentStatusPending:
		t.Errorf("expected: %v, actual: %v",
			PaymentInstallmentStatusPending,
			installment.Status,
		)
	}

	if dueAt, err := time.Parse(common.TimeFormat, installment.DueAt); err != nil || dueAt.IsZero() {
		t.Errorf("unexpected %v", installment.DueAt)
	}

	if id, err := uuid.Parse(installment.ID); err != nil || id == uuid.Nil {
		t.Errorf("unexpected %v", installment.ID)
	}
}

func TestPaymentServiceImp_CompletePaymentPlanCreation(t *testing.T) {
	t.Parallel()

	var (
		userID    = uuid.New()
		ctx       = context.Background()
		paymentID = uuid.New()
		dueAt, _  = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          paymentID,
		Currency:    "usdc",
		TotalAmount: "1000",
		Installments: []PaymentPlanInstallmentParams{
			{
				Currency: "usdc",
				Amount:   "1000",
				DueAt:    dueAt,
			},
		},
	}

	service := NewPaymentPlanService()

	_, err := service.CreatePendingPaymentPlan(ctx, userID, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	err = service.CompletePaymentPlanCreation(ctx, userID, paymentID)
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	payments, err := service.GetPaymentPlanByUserID(ctx, userID)
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	switch {
	case payments[0].ID != paymentPlanParams.ID.String():
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.ID, payments[0].ID)
	case payments[0].TotalAmount != paymentPlanParams.TotalAmount:
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.TotalAmount, payments[0].TotalAmount)
	case payments[0].Currency != paymentPlanParams.Currency:
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.Currency, payments[0].Currency)
	case payments[0].Status != paymentPlanStatusPending:
		t.Errorf("expected: %v, actual: %v", paymentPlanStatusPending, payments[0].Status)
	}

	installment := payments[0].Installments[0]
	installmentParams := paymentPlanParams.Installments[0]

	switch {
	case installment.Currency != installmentParams.Currency:
		t.Errorf("expected: %v, actual: %v",
			installment.Currency,
			installmentParams.Currency,
		)
	case installment.Amount != installmentParams.Amount:
		t.Errorf("expected: %v, actual: %v",
			installmentParams.Amount,
			installment.Amount,
		)
	case installment.Status != PaymentInstallmentStatusPaid:
		t.Errorf("expected: %v, actual: %v",
			PaymentInstallmentStatusPending,
			installment.Status,
		)
	}

	if dueAt, err := time.Parse(common.TimeFormat, installment.DueAt); err != nil || dueAt.IsZero() {
		t.Errorf("unexpected %v", installment.DueAt)
	}

	if id, err := uuid.Parse(installment.ID); err != nil || id == uuid.Nil {
		t.Errorf("unexpected %v", installment.ID)
	}
}
