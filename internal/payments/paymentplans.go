package payments

import (
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
)

type Plan struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Currency  string
	Amount    decimal.Big
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePlanParams struct {
	UserID   uuid.UUID
	Currency string
	Amount   decimal.Big
	Status   string
}
