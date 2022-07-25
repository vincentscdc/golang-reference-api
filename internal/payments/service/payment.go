package service

import (
	"context"
	"sort"

	"golangreferenceapi/internal/payments"
	"golangreferenceapi/internal/payments/common"
	"golangreferenceapi/internal/payments/repo"

	"github.com/ericlagergren/decimal"
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

var _ PaymentPlanService = (*PaymentServiceImp)(nil)

type PaymentServiceImp struct {
	repository repo.Repository
}

func NewPaymentPlanService() *PaymentServiceImp {
	return &PaymentServiceImp{}
}

func (p *PaymentServiceImp) UseRepo(repository repo.Repository) {
	p.repository = repository
}

func (p *PaymentServiceImp) GetPaymentPlanByUserID(ctx context.Context, userID uuid.UUID) ([]PaymentPlans, error) {
	plans, err := p.repository.ListPaymentPlansByUserID(ctx, userID)
	if err != nil {
		return nil, ListPaymentPlansByUserIDError{userID: userID}
	}

	paymentPlans := make([]PaymentPlans, 0, len(plans))

	for _, plan := range plans {
		paymentPlan := PaymentPlans{
			ID:          plan.ID.String(),
			UserID:      plan.UserID.String(),
			Currency:    plan.Currency,
			TotalAmount: plan.Amount.String(),
			Status:      plan.Status,
			CreatedAt:   plan.CreatedAt.Format(common.TimeFormat),
		}

		installments, err := p.repository.ListPaymentInstallmentsByPlanID(ctx, plan.ID)
		if err != nil {
			return nil, ListPaymentInstallmentsByPlanIDError{planID: plan.ID}
		}

		planInstallments := make([]PaymentPlanInstallment, 0, len(installments))

		for _, inst := range installments {
			planInst := PaymentPlanInstallment{
				ID:       inst.ID.String(),
				Amount:   inst.Amount.String(),
				Currency: inst.Currency,
				DueAt:    inst.DueAt.Format(common.TimeFormat),
				Status:   inst.Status,
			}
			planInstallments = append(planInstallments, planInst)
		}

		paymentPlan.Installments = planInstallments

		paymentPlans = append(paymentPlans, paymentPlan)
	}

	return paymentPlans, nil
}

func (p *PaymentServiceImp) CreatePendingPaymentPlan(
	ctx context.Context,
	paymentPlan *CreatePaymentPlanParams,
) (*PaymentPlans, error) {
	totalAmount := decimal.Big{}
	totalAmount.SetString(paymentPlan.TotalAmount)

	plan, err := p.repository.CreatePaymentPlan(ctx, &payments.CreatePlanParams{
		UserID:   paymentPlan.UserID,
		Currency: paymentPlan.Currency,
		Amount:   totalAmount,
		Status:   paymentPlanStatusPending,
	})
	if err != nil {
		return nil, CreatePaymentPlanError{}
	}

	newPlan := &PaymentPlans{
		ID:          plan.ID.String(),
		UserID:      plan.UserID.String(),
		Currency:    plan.Currency,
		TotalAmount: plan.Amount.String(),
		Status:      plan.Status,
		CreatedAt:   plan.CreatedAt.Format(common.TimeFormat),
	}

	sort.SliceStable(paymentPlan.Installments, func(i, j int) bool {
		return paymentPlan.Installments[i].DueAt.Unix() < paymentPlan.Installments[j].DueAt.Unix()
	})

	for _, inst := range paymentPlan.Installments {
		amount := decimal.Big{}
		amount.SetString(inst.Amount)

		installment, err := p.repository.CreatePaymentInstallment(ctx, &payments.CreateInstallmentParams{
			PaymentPlanID: plan.ID,
			Currency:      inst.Currency,
			Amount:        amount,
			DueAt:         inst.DueAt,
			Status:        PaymentInstallmentStatusPending,
		})
		if err != nil {
			return nil, CreatePaymentInstallmentError{}
		}

		newInst := PaymentPlanInstallment{
			ID:       installment.ID.String(),
			Amount:   installment.Amount.String(),
			Currency: installment.Currency,
			DueAt:    installment.DueAt.Format(common.TimeFormat),
			Status:   installment.Status,
		}

		newPlan.Installments = append(newPlan.Installments, newInst)
	}

	return newPlan, nil
}

// CompletePaymentPlanCreation Complete and paid the record of the first installments
func (p *PaymentServiceImp) CompletePaymentPlanCreation(
	ctx context.Context,
	paymentPlanID uuid.UUID,
	paymentPlan *CompletePaymentPlanParams,
) (*PaymentPlans, error) {
	var retPlan *PaymentPlans

	plans, err := p.repository.ListPaymentPlansByUserID(ctx, paymentPlan.UserID)
	if err != nil {
		return nil, ListPaymentPlansByUserIDError{userID: paymentPlan.UserID}
	}

	for _, plan := range plans {
		if plan.ID != paymentPlanID {
			continue
		}

		paymentPlan := PaymentPlans{
			ID:          plan.ID.String(),
			UserID:      plan.UserID.String(),
			Currency:    plan.Currency,
			TotalAmount: plan.Amount.String(),
			Status:      plan.Status,
			CreatedAt:   plan.CreatedAt.Format(common.TimeFormat),
		}

		installments, err := p.repository.ListPaymentInstallmentsByPlanID(ctx, plan.ID)
		if err != nil {
			return nil, ListPaymentInstallmentsByPlanIDError{planID: plan.ID}
		}

		planInstallments := make([]PaymentPlanInstallment, 0, len(installments))

		for _, inst := range installments {
			planInst := PaymentPlanInstallment{
				ID:       inst.ID.String(),
				Amount:   inst.Amount.String(),
				Currency: inst.Currency,
				DueAt:    inst.DueAt.Format(common.TimeFormat),
				Status:   inst.Status,
			}

			planInstallments = append(planInstallments, planInst)
		}

		paymentPlan.Installments = planInstallments

		retPlan = &paymentPlan

		return retPlan, nil
	}

	return nil, PaymentRecordNotFoundError{planID: paymentPlanID}
}
