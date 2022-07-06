package payments

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type Plan struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  string
	Amount    decimal.Decimal
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePlanParams struct {
	UserID   uuid.UUID
	Currency string
	Amount   decimal.Decimal
	Status   string
}
