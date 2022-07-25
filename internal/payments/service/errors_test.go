package service

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
)

func TestCreatePaymentPlanError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            CreatePaymentPlanError{},
			expectedString: "failed to create payment plan",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestListPaymentPlansByUserIDError(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ListPaymentPlansByUserIDError{userID: userID},
			expectedString: fmt.Sprintf("failed to get payment plans for user: %v", userID),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestCreatePaymentInstallmentError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            CreatePaymentInstallmentError{},
			expectedString: "failed to create payment plan installment",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestListPaymentInstallmentsByPlanIDError(t *testing.T) {
	t.Parallel()

	planID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            ListPaymentInstallmentsByPlanIDError{planID: planID},
			expectedString: fmt.Sprintf("failed to get payment plans for user: %v", planID),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}

func TestPaymentRecordNotFoundError(t *testing.T) {
	t.Parallel()

	planID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            PaymentRecordNotFoundError{planID: planID},
			expectedString: fmt.Sprintf("failed to get payment plan: %v", planID),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err.Error() != tt.expectedString {
				t.Error("unexpected Error string")
			}
		})
	}
}
