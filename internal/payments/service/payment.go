package service

import (
	"context"
	"sort"
	"time"

	"golangreferenceapi/internal/payments/common"
	"golangreferenceapi/internal/payments/repo"

	"github.com/gofrs/uuid"
)

const (
	paymentPlanStatusPending = "pending"
)

const (
	PaymentInstallmentStatusPending = "pending"
	PaymentInstallmentStatusPaid    = "paid"
	PaymentInstallmentStatusDue     = "due"
)

type UUIDGenerator func() (uuid.UUID, error)

var _ PaymentPlanService = (*PaymentServiceImp)(nil)

type PaymentServiceImp struct {
	uuidGenerator UUIDGenerator
	memoryStorage map[uuid.UUID][]PaymentPlans
	repository    repo.Repository
}

func NewPaymentPlanService() *PaymentServiceImp {
	return &PaymentServiceImp{
		uuidGenerator: uuid.NewV4,
		memoryStorage: make(map[uuid.UUID][]PaymentPlans),
	}
}

func (p *PaymentServiceImp) UseRepo(repository repo.Repository) {
	p.repository = repository
}

func (p *PaymentServiceImp) SetUUIDGenerator(generator UUIDGenerator) {
	p.uuidGenerator = generator
}

func (p *PaymentServiceImp) GetPaymentPlanByUserID(ctx context.Context, userID uuid.UUID) ([]PaymentPlans, error) {
	plans, ok := p.memoryStorage[userID]
	if !ok {
		return nil, ErrRecordNotFound
	}

	return plans, nil
}

func (p *PaymentServiceImp) CreatePendingPaymentPlan(
	ctx context.Context,
	paymentPlan *CreatePaymentPlanParams,
) (*PaymentPlans, error) {
	plan := &PaymentPlans{
		ID:          paymentPlan.ID.String(),
		UserID:      paymentPlan.UserID.String(),
		Currency:    paymentPlan.Currency,
		TotalAmount: paymentPlan.TotalAmount,
		Status:      paymentPlanStatusPending,
		CreatedAt:   time.Now().UTC().Format(common.TimeFormat),
	}

	sort.SliceStable(paymentPlan.Installments, func(i, j int) bool {
		return paymentPlan.Installments[i].DueAt.Unix() < paymentPlan.Installments[j].DueAt.Unix()
	})

	for _, inst := range paymentPlan.Installments {
		instID, err := p.uuidGenerator()
		if err != nil {
			return nil, ErrGenerateUUID
		}

		newInst := PaymentPlanInstallment{
			ID:       instID.String(),
			Currency: inst.Currency,
			Amount:   inst.Amount,
			DueAt:    inst.DueAt.Format(common.TimeFormat),
			Status:   PaymentInstallmentStatusPending,
		}

		plan.Installments = append(plan.Installments, newInst)
	}

	p.memoryStorage[paymentPlan.UserID] = append(p.memoryStorage[paymentPlan.UserID], *plan)

	return plan, nil
}

// CompletePaymentPlanCreation Complete and paid the record of the first installments
func (p *PaymentServiceImp) CompletePaymentPlanCreation(
	ctx context.Context,
	paymentPlanID uuid.UUID,
	paymentPlan *CompletePaymentPlanParams,
) (*PaymentPlans, error) {
	plan, ok := p.memoryStorage[paymentPlan.UserID]
	if !ok {
		return nil, ErrRecordNotFound
	}

	for i := range plan {
		if plan[i].ID == paymentPlanID.String() {
			plan[i].Installments[0].Status = PaymentInstallmentStatusPaid

			return &plan[i], nil
		}
	}

	return nil, ErrRecordNotFound
}
