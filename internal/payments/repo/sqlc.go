package repo

import (
	"context"

	"golangreferenceapi/internal/db"
	"golangreferenceapi/internal/payments"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SQLCRepo struct {
	querier db.Querier
}

func NewSQLCRepository(querier db.Querier) *SQLCRepo {
	return &SQLCRepo{querier: querier}
}

func (impl *SQLCRepo) CreatePaymentPlan(ctx context.Context, arg *payments.CreatePlanParams) (*payments.Plan, error) {
	amt, err := decimal.NewFromString(arg.Amount)
	if err != nil {
		return nil, decimalParseError{value: arg.Amount, WrappedErr: err}
	}

	dbEntity, err := impl.querier.CreatePaymentPlan(ctx, &db.CreatePaymentPlanParams{
		ID:       uuid.New(),
		UserID:   arg.UserID,
		Currency: db.Currency(arg.Currency),
		Amount:   amt,
		Status:   db.PaymentStatus(arg.Status),
	})
	if err != nil {
		return nil, pgxDBQueryRunError{WrappedErr: err}
	}

	plan, err := impl.newPlanFromDBEntity(dbEntity)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (impl *SQLCRepo) ListPaymentPlansByUserID(ctx context.Context, userID string) ([]*payments.Plan, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, uuidParseError{value: userID, WrappedErr: err}
	}

	entities, err := impl.querier.ListPaymentPlansByUserID(ctx, userUUID)
	if err != nil {
		return nil, pgxDBQueryRunError{WrappedErr: err}
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

func (impl *SQLCRepo) newPlanFromDBEntity(entity interface{}) (*payments.Plan, error) {
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

	return nil, unsupportedDBEntityError{}
}
