// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"fmt"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gofrs/uuid"
)

type Currency string

const (
	CurrencyUsdc Currency = "usdc"
)

func (e *Currency) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Currency(s)
	case string:
		*e = Currency(s)
	default:
		return fmt.Errorf("unsupported scan type for Currency: %T", src)
	}
	return nil
}

func (e Currency) Valid() bool {
	switch e {
	case CurrencyUsdc:
		return true
	}
	return false
}

func AllCurrencyValues() []Currency {
	return []Currency{
		CurrencyUsdc,
	}
}

type PaymentInstallmentStatus string

const (
	PaymentInstallmentStatusPending PaymentInstallmentStatus = "pending"
	PaymentInstallmentStatusPaid    PaymentInstallmentStatus = "paid"
	PaymentInstallmentStatusDue     PaymentInstallmentStatus = "due"
)

func (e *PaymentInstallmentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentInstallmentStatus(s)
	case string:
		*e = PaymentInstallmentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentInstallmentStatus: %T", src)
	}
	return nil
}

func (e PaymentInstallmentStatus) Valid() bool {
	switch e {
	case PaymentInstallmentStatusPending,
		PaymentInstallmentStatusPaid,
		PaymentInstallmentStatusDue:
		return true
	}
	return false
}

func AllPaymentInstallmentStatusValues() []PaymentInstallmentStatus {
	return []PaymentInstallmentStatus{
		PaymentInstallmentStatusPending,
		PaymentInstallmentStatusPaid,
		PaymentInstallmentStatusDue,
	}
}

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusComplete PaymentStatus = "complete"
)

func (e *PaymentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatus(s)
	case string:
		*e = PaymentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatus: %T", src)
	}
	return nil
}

func (e PaymentStatus) Valid() bool {
	switch e {
	case PaymentStatusPending,
		PaymentStatusComplete:
		return true
	}
	return false
}

func AllPaymentStatusValues() []PaymentStatus {
	return []PaymentStatus{
		PaymentStatusPending,
		PaymentStatusComplete,
	}
}

type PaymentInstallment struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Currency      Currency
	Amount        decimal.Big
	DueAt         time.Time
	Status        PaymentInstallmentStatus
	PaymentPlanID uuid.UUID
}

type PaymentPlan struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Currency  Currency
	UserID    uuid.UUID
	Amount    decimal.Big
	Status    PaymentStatus
}
