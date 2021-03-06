// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: payment_plans.sql

package db

import (
	"context"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
)

const CreatePaymentPlan = `-- name: CreatePaymentPlan :one
INSERT INTO payment_plans (id, user_id, currency, amount, status) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, currency, amount, status, created_at, updated_at
`

type CreatePaymentPlanParams struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Currency Currency
	Amount   decimal.Big
	Status   PaymentStatus
}

type CreatePaymentPlanRow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  Currency
	Amount    decimal.Big
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreatePaymentPlan(ctx context.Context, arg *CreatePaymentPlanParams) (*CreatePaymentPlanRow, error) {
	row := q.db.QueryRow(ctx, CreatePaymentPlan,
		arg.ID,
		arg.UserID,
		arg.Currency,
		arg.Amount,
		arg.Status,
	)
	var i CreatePaymentPlanRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Currency,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const ListPaymentPlansByUserID = `-- name: ListPaymentPlansByUserID :many
SELECT id, user_id, currency, amount, status, created_at, updated_at FROM payment_plans
WHERE user_id = $1
ORDER BY created_at DESC
`

type ListPaymentPlansByUserIDRow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  Currency
	Amount    decimal.Big
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) ListPaymentPlansByUserID(ctx context.Context, userID uuid.UUID) ([]*ListPaymentPlansByUserIDRow, error) {
	rows, err := q.db.Query(ctx, ListPaymentPlansByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListPaymentPlansByUserIDRow
	for rows.Next() {
		var i ListPaymentPlansByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Currency,
			&i.Amount,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
