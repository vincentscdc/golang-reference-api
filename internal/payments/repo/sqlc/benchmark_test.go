package sqlc

import (
	"context"
	"testing"
	"time"

	"golangreferenceapi/internal/payments"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// nolint: gochecknoglobals // tests
var (
	testBenchCtx                    = context.Background()
	testBenchAmount, _              = decimal.NewFromString("10.98")
	testBenchUserID                 = uuid.Must(uuid.NewV4())
	testBenchCreatePaymentPlanParam = &payments.CreatePlanParams{
		UserID:   testBenchUserID,
		Currency: "usdc",
		Amount:   testBenchAmount,
		Status:   "pending",
	}
	testBenchRefPaymentPlan = &payments.Plan{
		ID:        uuid.UUID{},
		UserID:    uuid.UUID{},
		Currency:  "",
		Amount:    decimal.Decimal{},
		Status:    "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
)

func BenchmarkCreatePaymentPlan(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		testBenchRefPaymentPlan, err = testRefRepo.CreatePaymentPlan(testBenchCtx, testBenchCreatePaymentPlanParam)
		if err != nil {
			b.Fatalf("err creating payment plan")
		}
	}
}

func BenchmarkListPaymentPlansByUserID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := testRefRepo.ListPaymentPlansByUserID(testBenchCtx, testBenchUserID)
		if err != nil {
			b.Fatalf("err listing payment plans by user id")
		}
	}
}

func BenchmarkCreatePaymentInstallment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := testRefRepo.CreatePaymentInstallment(testBenchCtx, &payments.CreateInstallmentParams{
			PaymentPlanID: testBenchRefPaymentPlan.ID,
			Currency:      "usdc",
			Amount:        testBenchAmount,
			DueAt:         time.Time{},
			Status:        "pending",
		})
		if err != nil {
			b.Fatalf("err creating payment installment")
		}
	}
}

func BenchmarkListPaymentInstallmentsByPlanID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := testRefRepo.ListPaymentInstallmentsByPlanID(testBenchCtx, testBenchRefPaymentPlan.ID)
		if err != nil {
			b.Fatalf("err listing payment installments by plan id")
		}
	}
}
