package rest

import (
	"errors"
	"net/http"
	"testing"

	"golangreferenceapi/internal/payments/service"
)

func TestServiceErrorToErrorResp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		statusCode int
	}{
		{
			name:       "create payment plan error",
			err:        service.CreatePaymentPlanError{},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "list payment plan by userid error",
			err:        service.ListPaymentPlansByUserIDError{},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "create payment installment error",
			err:        service.CreatePaymentInstallmentError{},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "list payment installment by planid error",
			err:        service.ListPaymentInstallmentsByPlanIDError{},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "payment record not found",
			err:        service.PaymentRecordNotFoundError{},
			statusCode: http.StatusNotFound,
		},
		{
			name:       "unknown",
			err:        errors.New("error unknown"),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errResp := ServiceErrorToErrorResp(tt.err)
			if errResp.StatusCode != tt.statusCode {
				t.Errorf("unexpected status code, expected: %v, actual: %v",
					tt.statusCode,
					errResp.StatusCode,
				)
			}
		})
	}
}
