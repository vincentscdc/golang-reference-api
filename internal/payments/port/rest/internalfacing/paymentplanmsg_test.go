// nolint: cyclop // disable for testing
package internalfacing

import (
	"reflect"
	"testing"
	"time"

	"golangreferenceapi/internal/payments/common"

	"github.com/google/uuid"
)

func TestNewCreatePaymentPlanParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		paymentID        string
		installmentID    string
		installmentDueAt string
		expectedErr      error
	}{
		{
			name:             "Happy Path",
			paymentID:        uuid.New().String(),
			installmentID:    uuid.New().String(),
			installmentDueAt: time.Now().Format(common.TimeFormat),
			expectedErr:      nil,
		},
		{
			name:             "invalid payment id",
			paymentID:        "not parsable id",
			installmentID:    uuid.New().String(),
			installmentDueAt: time.Now().Format(common.TimeFormat),
			expectedErr:      inputParseError{},
		},
		{
			name:             "invalid installment id",
			paymentID:        uuid.New().String(),
			installmentID:    "not parsable id",
			installmentDueAt: time.Now().Format(common.TimeFormat),
			expectedErr:      inputParseError{},
		},
		{
			name:             "wrong time format",
			paymentID:        uuid.New().String(),
			installmentID:    uuid.New().String(),
			installmentDueAt: time.ANSIC,
			expectedErr:      inputParseError{},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := &CreatePendingPaymentPlanRequest{
				PendingPayment: CreatePendingPayment{
					ID:          tt.paymentID,
					Currency:    "usdc",
					TotalAmount: "200",
					Installments: []CreateInstallment{
						{
							ID:    tt.installmentID,
							DueAt: tt.installmentDueAt,
						},
					},
				},
			}
			params, err := NewCreatePaymentPlanParams(req)
			if (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr == nil) {
				t.Fatalf(
					"unknown error expected: %v, actual: %v",
					tt.expectedErr,
					err,
				)
			} else if reflect.TypeOf(err) != reflect.TypeOf(tt.expectedErr) {
				t.Fatalf(
					"unknown error expected: %v, actual: %v",
					reflect.TypeOf(err).Name(),
					reflect.TypeOf(tt.expectedErr).Name(),
				)
			} else if err != nil {
				return
			}

			if params.ID.String() != tt.paymentID {
				t.Errorf("expected %v actual %v", tt.paymentID, params.ID.String())
			}

			if params.Currency != req.PendingPayment.Currency {
				t.Errorf("expected %v actual %v", req.PendingPayment.Currency, params.ID.String())
			}

			if params.TotalAmount != req.PendingPayment.TotalAmount {
				t.Errorf("expected %v actual %v", req.PendingPayment.TotalAmount, params.TotalAmount)
			}

			if params.Installments[0].ID.String() != req.PendingPayment.Installments[0].ID {
				t.Errorf(
					"expected %v actual %v",
					req.PendingPayment.Installments[0].ID,
					params.Installments[0].ID.String(),
				)
			}

			if params.Installments[0].DueAt.Format(common.TimeFormat) != req.PendingPayment.Installments[0].DueAt {
				t.Errorf("expected %v actual %v",
					req.PendingPayment.Installments[0].DueAt,
					params.Installments[0].DueAt.Format(common.TimeFormat),
				)
			}
		})
	}
}
