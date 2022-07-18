package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"golangreferenceapi/internal/payments/common"
	"golangreferenceapi/internal/payments/repo/sqlc"

	"github.com/gofrs/uuid"
)

func TestPaymentServiceImp_CreatePendingPaymentPlan(t *testing.T) {
	t.Parallel()

	var (
		userID        = uuid.Must(uuid.NewV4())
		paymentPlanID = uuid.Must(uuid.NewV4())
		ctx           = context.Background()
		dueAt, _      = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          paymentPlanID,
		UserID:      userID,
		Currency:    "usdc",
		TotalAmount: "1000",
		Installments: []PaymentPlanInstallmentParams{
			{
				Currency: "usdc",
				Amount:   "1000",
				DueAt:    dueAt,
			},
			{
				Currency: "usdc",
				Amount:   "1000",
				DueAt:    dueAt.Add(1 * time.Hour),
			},
			{
				Currency: "usdc",
				Amount:   "1000",
				DueAt:    dueAt.Add(2 * time.Hour),
			},
		},
	}

	service := NewPaymentPlanService()

	_, err := service.CreatePendingPaymentPlan(ctx, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	payments, err := service.GetPaymentPlanByUserID(ctx, userID)
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	if payments[0].ID != paymentPlanParams.ID.String() {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.ID, payments[0].ID)
	}

	if payments[0].TotalAmount != paymentPlanParams.TotalAmount {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.TotalAmount, payments[0].TotalAmount)
	}

	if payments[0].Currency != paymentPlanParams.Currency {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.Currency, payments[0].Currency)
	}

	if payments[0].Status != paymentPlanStatusPending {
		t.Errorf("expected: %v, actual: %v", paymentPlanStatusPending, payments[0].Status)
	}

	installment := payments[0].Installments[0]
	installmentParams := paymentPlanParams.Installments[0]

	if installment.Currency != installmentParams.Currency {
		t.Errorf("expected: %v, actual: %v",
			installment.Currency,
			installmentParams.Currency,
		)
	}

	if installment.Amount != installmentParams.Amount {
		t.Errorf("expected: %v, actual: %v",
			installmentParams.Amount,
			installment.Amount,
		)
	}

	if installment.Status != PaymentInstallmentStatusPending {
		t.Errorf("expected: %v, actual: %v",
			PaymentInstallmentStatusPending,
			installment.Status,
		)
	}

	if dueAt, err := time.Parse(common.TimeFormat, installment.DueAt); err != nil || dueAt.IsZero() {
		t.Errorf("unexpected %v", installment.DueAt)
	}

	if id, err := uuid.FromString(installment.ID); err != nil || id == uuid.Nil {
		t.Errorf("unexpected %v", installment.ID)
	}
}

func TestPaymentServiceImp_GetPaymentPlanByUserID(t *testing.T) {
	t.Parallel()

	var (
		userID   = uuid.Must(uuid.NewV4())
		planID   = uuid.Must(uuid.NewV4())
		ctx      = context.Background()
		dueAt, _ = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          planID,
		UserID:      userID,
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

	_, err := service.CreatePendingPaymentPlan(ctx, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	payments, err := service.GetPaymentPlanByUserID(ctx, uuid.Must(uuid.NewV4()))
	if !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("unexpected record existed: %v", payments)
	}

	payments, err = service.GetPaymentPlanByUserID(ctx, userID)
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	if payments[0].ID != paymentPlanParams.ID.String() {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.ID, payments[0].ID)
	}

	if payments[0].TotalAmount != paymentPlanParams.TotalAmount {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.TotalAmount, payments[0].TotalAmount)
	}

	if payments[0].Currency != paymentPlanParams.Currency {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.Currency, payments[0].Currency)
	}

	if payments[0].Status != paymentPlanStatusPending {
		t.Errorf("expected: %v, actual: %v", paymentPlanStatusPending, payments[0].Status)
	}

	installment := payments[0].Installments[0]
	installmentParams := paymentPlanParams.Installments[0]

	if installment.Currency != installmentParams.Currency {
		t.Errorf("expected: %v, actual: %v",
			installment.Currency,
			installmentParams.Currency,
		)
	}

	if installment.Amount != installmentParams.Amount {
		t.Errorf("expected: %v, actual: %v",
			installmentParams.Amount,
			installment.Amount,
		)
	}

	if installment.Status != PaymentInstallmentStatusPending {
		t.Errorf("expected: %v, actual: %v",
			PaymentInstallmentStatusPending,
			installment.Status,
		)
	}

	if dueAt, err := time.Parse(common.TimeFormat, installment.DueAt); err != nil || dueAt.IsZero() {
		t.Errorf("unexpected %v", installment.DueAt)
	}

	if id, err := uuid.FromString(installment.ID); err != nil || id == uuid.Nil {
		t.Errorf("unexpected %v", installment.ID)
	}
}

func TestPaymentServiceImp_CompletePaymentPlanCreation(t *testing.T) {
	t.Parallel()

	var (
		userID    = uuid.Must(uuid.NewV4())
		ctx       = context.Background()
		paymentID = uuid.Must(uuid.NewV4())
		dueAt, _  = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          paymentID,
		UserID:      userID,
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

	_, err := service.CreatePendingPaymentPlan(ctx, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	payment, err := service.CompletePaymentPlanCreation(ctx, paymentID, &CompletePaymentPlanParams{
		UserID: userID,
	})
	if err != nil {
		t.Errorf("failed to get payment plan %s", err.Error())
	}

	if payment.ID != paymentPlanParams.ID.String() {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.ID, payment.ID)
	}

	if payment.TotalAmount != paymentPlanParams.TotalAmount {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.TotalAmount, payment.TotalAmount)
	}

	if payment.Currency != paymentPlanParams.Currency {
		t.Errorf("expected: %v, actual: %v", paymentPlanParams.Currency, payment.Currency)
	}

	if payment.Status != paymentPlanStatusPending {
		t.Errorf("expected: %v, actual: %v", paymentPlanStatusPending, payment.Status)
	}

	installment := payment.Installments[0]
	installmentParams := paymentPlanParams.Installments[0]

	if installment.Currency != installmentParams.Currency {
		t.Errorf("expected: %v, actual: %v",
			installment.Currency,
			installmentParams.Currency,
		)
	}

	if installment.Amount != installmentParams.Amount {
		t.Errorf("expected: %v, actual: %v",
			installmentParams.Amount,
			installment.Amount,
		)
	}

	if installment.Status != PaymentInstallmentStatusPaid {
		t.Errorf("expected: %v, actual: %v",
			PaymentInstallmentStatusPending,
			installment.Status,
		)
	}

	if dueAt, err := time.Parse(common.TimeFormat, installment.DueAt); err != nil || dueAt.IsZero() {
		t.Errorf("unexpected %v", installment.DueAt)
	}

	if id, err := uuid.FromString(installment.ID); err != nil || id == uuid.Nil {
		t.Errorf("unexpected %v", installment.ID)
	}
}

func TestPaymentServiceImp_CompletePaymentPlanCreationNotFound(t *testing.T) {
	t.Parallel()

	var (
		userID    = uuid.Must(uuid.NewV4())
		ctx       = context.Background()
		paymentID = uuid.Must(uuid.NewV4())
		dueAt, _  = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	tests := []struct {
		name      string
		userID    uuid.UUID
		paymentID uuid.UUID
	}{
		{
			name:      "payment record not found",
			userID:    userID,
			paymentID: uuid.Must(uuid.NewV4()),
		},
		{
			name:      "user record not fund",
			userID:    uuid.Must(uuid.NewV4()),
			paymentID: paymentID,
		},
	}

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          paymentID,
		UserID:      userID,
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

	_, err := service.CreatePendingPaymentPlan(ctx, paymentPlanParams)
	if err != nil {
		t.Errorf("failed to create payment plan %s", err.Error())
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.CompletePaymentPlanCreation(ctx, tt.paymentID, &CompletePaymentPlanParams{
				UserID: tt.userID,
			})
			if !errors.Is(err, ErrRecordNotFound) {
				t.Errorf("failed to get error, expected: %v, got: %v", ErrRecordNotFound, err)
			}
		})
	}
}

func TestPaymentServiceImp_CompletePaymentPlanCreationUUIDGenError(t *testing.T) {
	t.Parallel()

	var (
		userID    = uuid.Must(uuid.NewV4())
		ctx       = context.Background()
		paymentID = uuid.Must(uuid.NewV4())
		dueAt, _  = time.Parse(common.TimeFormat, "2021-11-10T23:00:00Z")
	)

	paymentPlanParams := &CreatePaymentPlanParams{
		ID:          paymentID,
		UserID:      userID,
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

	failedUUIDGen := func() (uuid.UUID, error) {
		return uuid.Nil, errors.New("failed to generate uuid")
	}

	service := NewPaymentPlanService()
	service.SetUUIDGenerator(failedUUIDGen)

	_, err := service.CreatePendingPaymentPlan(ctx, paymentPlanParams)
	if !errors.Is(err, ErrGenerateUUID) {
		t.Errorf("error expected: %v, actual: %v", ErrGenerateUUID, err)
	}
}

func TestPaymentServiceImp_UseRepo(t *testing.T) {
	t.Parallel()

	service := NewPaymentPlanService()
	service.UseRepo(&sqlc.Repo{})

	if reflect.TypeOf(service.repository) != reflect.TypeOf(&sqlc.Repo{}) {
		t.Errorf("returned repository is not of Repo")
	}
}
