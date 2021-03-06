// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: payment_installments.sql

package db

import (
	"context"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
)

const CreatePaymentInstallments = `-- name: CreatePaymentInstallments :one
INSERT INTO payment_installments (id, payment_plan_id, currency, amount, due_at, status) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, payment_plan_id, currency, amount, due_at, status, created_at, updated_at
`

type CreatePaymentInstallmentsParams struct {
	ID            uuid.UUID
	PaymentPlanID uuid.UUID
	Currency      Currency
	Amount        decimal.Big
	DueAt         time.Time
	Status        PaymentInstallmentStatus
}

type CreatePaymentInstallmentsRow struct {
	ID            uuid.UUID
	PaymentPlanID uuid.UUID
	Currency      Currency
	Amount        decimal.Big
	DueAt         time.Time
	Status        PaymentInstallmentStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (q *Queries) CreatePaymentInstallments(ctx context.Context, arg *CreatePaymentInstallmentsParams) (*CreatePaymentInstallmentsRow, error) {
	row := q.db.QueryRow(ctx, CreatePaymentInstallments,
		arg.ID,
		arg.PaymentPlanID,
		arg.Currency,
		arg.Amount,
		arg.DueAt,
		arg.Status,
	)
	var i CreatePaymentInstallmentsRow
	err := row.Scan(
		&i.ID,
		&i.PaymentPlanID,
		&i.Currency,
		&i.Amount,
		&i.DueAt,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const ListPaymentInstallmentsByPlanID = `-- name: ListPaymentInstallmentsByPlanID :many
SELECT id, payment_plan_id, currency, amount, due_at, status, created_at, updated_at FROM payment_installments
WHERE payment_plan_id = $1
ORDER BY due_at
`

type ListPaymentInstallmentsByPlanIDRow struct {
	ID            uuid.UUID
	PaymentPlanID uuid.UUID
	Currency      Currency
	Amount        decimal.Big
	DueAt         time.Time
	Status        PaymentInstallmentStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (q *Queries) ListPaymentInstallmentsByPlanID(ctx context.Context, paymentPlanID uuid.UUID) ([]*ListPaymentInstallmentsByPlanIDRow, error) {
	rows, err := q.db.Query(ctx, ListPaymentInstallmentsByPlanID, paymentPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListPaymentInstallmentsByPlanIDRow
	for rows.Next() {
		var i ListPaymentInstallmentsByPlanIDRow
		if err := rows.Scan(
			&i.ID,
			&i.PaymentPlanID,
			&i.Currency,
			&i.Amount,
			&i.DueAt,
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
