package sqlc

import (
	"context"

	"golangreferenceapi/internal/db"
	"golangreferenceapi/internal/payments"

	"github.com/google/uuid"
)

type Repo struct {
	querier db.Querier
}

func NewSQLCRepository(querier db.Querier) *Repo {
	return &Repo{querier: querier}
}

func (impl *Repo) CreatePaymentPlan(ctx context.Context, arg *payments.CreatePlanParams) (*payments.Plan, error) {
	dbEntity, err := impl.querier.CreatePaymentPlan(ctx, &db.CreatePaymentPlanParams{
		ID:       uuid.New(),
		UserID:   arg.UserID,
		Currency: db.Currency(arg.Currency),
		Amount:   arg.Amount,
		Status:   db.PaymentStatus(arg.Status),
	})
	if err != nil {
		return nil, err
	}

	plan, err := impl.newPlanFromDBEntity(dbEntity)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (impl *Repo) ListPaymentPlansByUserID(ctx context.Context, userID uuid.UUID) ([]*payments.Plan, error) {
	entities, err := impl.querier.ListPaymentPlansByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	plans := make([]*payments.Plan, len(entities))

	for idx, entity := range entities {
		plan, err := impl.newPlanFromDBEntity(entity)
		if err != nil {
			return nil, err
		}

		plans[idx] = plan
	}

	return plans, nil
}

func (impl *Repo) CreatePaymentInstallment(
	ctx context.Context,
	arg *payments.CreateInstallmentParams,
) (*payments.Installment, error) {
	dbEntity, err := impl.querier.CreatePaymentInstallments(ctx, &db.CreatePaymentInstallmentsParams{
		ID:            uuid.New(),
		PaymentPlanID: arg.PaymentPlanID,
		Currency:      db.Currency(arg.Currency),
		Amount:        arg.Amount,
		DueAt:         arg.DueAt,
		Status:        db.PaymentInstallmentStatus(arg.Status),
	})
	if err != nil {
		return nil, err
	}

	installment, err := impl.newInstallmentFromDBEntity(dbEntity)
	if err != nil {
		return nil, err
	}

	return installment, nil
}

func (impl *Repo) ListPaymentInstallmentsByPlanID(
	ctx context.Context,
	planID uuid.UUID,
) ([]*payments.Installment, error) {
	entities, err := impl.querier.ListPaymentInstallmentsByPlanID(ctx, planID)
	if err != nil {
		return nil, err
	}

	installments := make([]*payments.Installment, len(entities))

	for idx, entity := range entities {
		plan, err := impl.newInstallmentFromDBEntity(entity)
		if err != nil {
			return nil, err
		}

		installments[idx] = plan
	}

	return installments, nil
}

func (impl *Repo) newPlanFromDBEntity(entity interface{}) (*payments.Plan, error) {
	createPaymentPlanRowEntity, valid := entity.(*db.CreatePaymentPlanRow)
	if valid {
		return &payments.Plan{
			ID:        createPaymentPlanRowEntity.ID,
			UserID:    createPaymentPlanRowEntity.UserID,
			Currency:  string(createPaymentPlanRowEntity.Currency),
			Amount:    createPaymentPlanRowEntity.Amount,
			Status:    string(createPaymentPlanRowEntity.Status),
			CreatedAt: createPaymentPlanRowEntity.CreatedAt,
			UpdatedAt: createPaymentPlanRowEntity.UpdatedAt,
		}, nil
	}

	listPaymentPlansByUserIDRowEntity, valid := entity.(*db.ListPaymentPlansByUserIDRow)
	if valid {
		return &payments.Plan{
			ID:        listPaymentPlansByUserIDRowEntity.ID,
			UserID:    listPaymentPlansByUserIDRowEntity.UserID,
			Currency:  string(listPaymentPlansByUserIDRowEntity.Currency),
			Amount:    listPaymentPlansByUserIDRowEntity.Amount,
			Status:    string(listPaymentPlansByUserIDRowEntity.Status),
			CreatedAt: listPaymentPlansByUserIDRowEntity.CreatedAt,
			UpdatedAt: listPaymentPlansByUserIDRowEntity.UpdatedAt,
		}, nil
	}

	planEntity, valid := entity.(*db.PaymentPlan)
	if valid {
		return &payments.Plan{
			ID:        planEntity.ID,
			UserID:    planEntity.UserID,
			Currency:  string(planEntity.Currency),
			Amount:    planEntity.Amount,
			Status:    string(planEntity.Status),
			CreatedAt: planEntity.CreatedAt,
			UpdatedAt: planEntity.UpdatedAt,
		}, nil
	}

	return nil, UnsupportedDBEntityError{}
}

func (impl *Repo) newInstallmentFromDBEntity(entity interface{}) (*payments.Installment, error) {
	createInstRowEntity, valid := entity.(*db.CreatePaymentInstallmentsRow)
	if valid {
		return &payments.Installment{
			ID:            createInstRowEntity.ID,
			PaymentPlanID: createInstRowEntity.PaymentPlanID,
			Currency:      string(createInstRowEntity.Currency),
			Amount:        createInstRowEntity.Amount,
			DueAt:         createInstRowEntity.DueAt,
			Status:        string(createInstRowEntity.Status),
			CreatedAt:     createInstRowEntity.CreatedAt,
			UpdatedAt:     createInstRowEntity.UpdatedAt,
		}, nil
	}

	listInstsByUserIDRowEntity, valid := entity.(*db.ListPaymentInstallmentsByPlanIDRow)
	if valid {
		return &payments.Installment{
			ID:            listInstsByUserIDRowEntity.ID,
			PaymentPlanID: listInstsByUserIDRowEntity.PaymentPlanID,
			Currency:      string(listInstsByUserIDRowEntity.Currency),
			Amount:        listInstsByUserIDRowEntity.Amount,
			DueAt:         listInstsByUserIDRowEntity.DueAt,
			Status:        string(listInstsByUserIDRowEntity.Status),
			CreatedAt:     listInstsByUserIDRowEntity.CreatedAt,
			UpdatedAt:     listInstsByUserIDRowEntity.UpdatedAt,
		}, nil
	}

	instEntity, valid := entity.(*db.PaymentInstallment)
	if valid {
		return &payments.Installment{
			ID:            instEntity.ID,
			PaymentPlanID: instEntity.PaymentPlanID,
			Currency:      string(instEntity.Currency),
			Amount:        instEntity.Amount,
			DueAt:         instEntity.DueAt,
			Status:        string(instEntity.Status),
			CreatedAt:     instEntity.CreatedAt,
			UpdatedAt:     instEntity.UpdatedAt,
		}, nil
	}

	return nil, UnsupportedDBEntityError{}
}
