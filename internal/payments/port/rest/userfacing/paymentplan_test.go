package userfacing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"golangreferenceapi/internal/payments/service"

	"github.com/google/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func Test_getPaymentPlansHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		query                 string
		expectedResponse      *handlerwrap.Response
		expectedErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name:                  "happy path",
			query:                 "offset=10&limit=10&created_at_order=desc",
			expectedErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusNotFound},
			expectedResponse: &handlerwrap.Response{
				Body:           ([]PaymentPlanResponse)(nil),
				HTTPStatusCode: http.StatusOK,
			},
		},
		{
			name:                  "invalid pagination query params",
			query:                 "offset=x&limit=x&created_at_order=x",
			expectedErrorResponse: handlerwrap.ParsingParamError{}.ToErrorResponse(),
		},
	}

	paymentService := service.NewPaymentPlanService()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				req := httptest.NewRequest("GET", "/?"+tt.query, nil)
				id, _ := uuid.Parse("b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5")
				req = req.WithContext(context.WithValue(req.Context(), contextValKeyUserUUID, &id))

				resp, err := listPaymentPlansHandler(paymentService)(req)
				if tt.expectedErrorResponse != nil {
					if err.HTTPStatusCode != tt.expectedErrorResponse.HTTPStatusCode {
						t.Errorf("returned a unexpected error code got %v want %v", err.HTTPStatusCode, tt.expectedErrorResponse.HTTPStatusCode)
					}

					return
				}
				if !reflect.DeepEqual(resp, tt.expectedResponse) {
					t.Errorf("returned a unexpected response got %#v want %#v", resp, tt.expectedResponse)
				}
			})
		})
	}
}

func Test_getPaymentPlansHandler_MissingUserID(t *testing.T) {
	t.Parallel()

	paymentService := service.NewPaymentPlanService()

	req := httptest.NewRequest("GET", "/", nil)

	_, err := listPaymentPlansHandler(paymentService)(req)
	if err != nil && err.HTTPStatusCode != http.StatusUnauthorized {
		t.Errorf("unexpected status code expect: %v, actual: %v", err.ErrorCode, http.StatusBadRequest)
	}
}

func Benchmark_getPaymentPlansHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	paymentService := service.NewPaymentPlanService()
	h := listPaymentPlansHandler(paymentService)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
