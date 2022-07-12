package payments

import (
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
)

type Installment struct {
	ID            uuid.UUID   `json:"id"`
	PaymentPlanID uuid.UUID   `json:"payment_plan_id"`
	Currency      string      `json:"currency"`
	Amount        decimal.Big `json:"amount"`
	DueAt         time.Time   `json:"due_at"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type CreateInstallmentParams struct {
	PaymentPlanID uuid.UUID
	Currency      string
	Amount        decimal.Big
	DueAt         time.Time
	Status        string
}
