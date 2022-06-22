package payments

import (
	"time"

	"github.com/shopspring/decimal"
)

type Installment struct {
	ID            string
	PaymentPlanID string
	Currency      string
	Amount        decimal.Decimal
	DueAt         time.Time
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
