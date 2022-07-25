package memory

import (
	"context"
	"sync"
	"time"

	"golangreferenceapi/internal/payments"

	"github.com/gofrs/uuid"
)

type InMemRepo struct {
	paymentPlansLock        sync.RWMutex
	paymentPlans            map[uuid.UUID][]*payments.Plan
	paymentInstallmentsLock sync.RWMutex
	paymentInstallments     map[uuid.UUID][]*payments.Installment
}

type memoryError string

func (me memoryError) Error() string {
	return string(me)
}

const (
	ErrRecordNotFound   = memoryError("no records found")
	ErrGenerateUUID     = memoryError("failed to generate uuid")
	ErrMapTypeAssertion = memoryError("type assertion failed when load map")
)

func NewInMemRepository() *InMemRepo {
	return &InMemRepo{
		paymentPlans:        make(map[uuid.UUID][]*payments.Plan),
		paymentInstallments: make(map[uuid.UUID][]*payments.Installment),
	}
}

func (imr *InMemRepo) CreatePaymentPlan(
	ctx context.Context,
	arg *payments.CreatePlanParams,
) (*payments.Plan, error) {
	planID, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID
	}

	plan := &payments.Plan{
		ID:        planID,
		UserID:    arg.UserID,
		Currency:  arg.Currency,
		Amount:    arg.Amount,
		Status:    arg.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	imr.paymentPlansLock.Lock()
	imr.paymentPlans[arg.UserID] = append(imr.paymentPlans[arg.UserID], plan)
	imr.paymentPlansLock.Unlock()

	return plan, nil
}

func (imr *InMemRepo) ListPaymentPlansByUserID(ctx context.Context, userID uuid.UUID) ([]*payments.Plan, error) {
	imr.paymentPlansLock.RLock()
	defer imr.paymentPlansLock.RUnlock()

	res, ok := imr.paymentPlans[userID]
	if !ok {
		return nil, ErrRecordNotFound
	}

	return res, nil
}

func (imr *InMemRepo) CreatePaymentInstallment(
	ctx context.Context,
	arg *payments.CreateInstallmentParams,
) (*payments.Installment, error) {
	installmentID, err := uuid.NewV4()
	if err != nil {
		return nil, ErrGenerateUUID
	}

	inst := &payments.Installment{
		ID:            installmentID,
		PaymentPlanID: arg.PaymentPlanID,
		Currency:      arg.Currency,
		Amount:        arg.Amount,
		DueAt:         arg.DueAt,
		Status:        arg.Status,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	imr.paymentInstallmentsLock.Lock()
	imr.paymentInstallments[arg.PaymentPlanID] = append(imr.paymentInstallments[arg.PaymentPlanID], inst)
	imr.paymentInstallmentsLock.Unlock()

	return inst, nil
}

func (imr *InMemRepo) ListPaymentInstallmentsByPlanID(
	ctx context.Context,
	planID uuid.UUID,
) ([]*payments.Installment, error) {
	imr.paymentInstallmentsLock.RLock()
	defer imr.paymentInstallmentsLock.RUnlock()

	res, ok := imr.paymentInstallments[planID]
	if !ok {
		return nil, ErrRecordNotFound
	}

	return res, nil
}
