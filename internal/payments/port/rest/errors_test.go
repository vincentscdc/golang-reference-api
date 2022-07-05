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
			name:       "uuid generated error",
			err:        service.ErrGenerateUUID,
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "record not found",
			err:        service.ErrRecordNotFound,
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
			if errResp.HTTPStatusCode != tt.statusCode {
				t.Errorf("unexpected status code, expected: %v, actual: %v",
					tt.statusCode,
					errResp.HTTPStatusCode,
				)
			}
		})
	}
}
