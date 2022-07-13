package internalfacing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"golangreferenceapi/internal/payments/mock/servicemock"
	"golangreferenceapi/internal/payments/service"
	"golangreferenceapi/internal/payments/transport/rest"
)

func TestAddRoutes(t *testing.T) {
	t.Parallel()

	log := zerolog.Nop().With().Logger()

	userID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name                   string
		httpMethod             string
		urlPath                string
		reqBody                string
		expectedHTTPStatusCode int
	}{
		{
			name:       "happy path for creating a pending payment plan",
			httpMethod: "POST",
			urlPath:    fmt.Sprintf("/api/internal/pay_later/users/%s/payment_plans", userID.String()),
			reqBody: `{
					"payment": {
						"id": "03baa9e6-6ed6-4868-9ef9-b99c8452f270",
						"currency": "usdc",
						"total_amount": "100.0",
						"installments": [
							{ "due_at": "2022-06-01T14:02:03.000Z", "amount": "50", "currency": "usdc"},
							{ "due_at": "2022-06-15T14:02:03.000Z", "amount": "50", "currency": "usdc"}
						]
				  }
				}`,
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:       "happy path for completing payment plan",
			httpMethod: "POST",
			urlPath: fmt.Sprintf(
				"/api/internal/pay_later/users/%s/payment_plans/03baa9e6-6ed6-4868-9ef9-b99c8452f270/complete",
				userID.String(),
			),
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
		},
	}

	paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))
	paymentService.EXPECT().
		CreatePendingPaymentPlan(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&service.PaymentPlans{}, nil)

	paymentService.EXPECT().
		CompletePaymentPlanCreation(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&service.PaymentPlans{}, nil)

	for _, tt := range tests { //nolint: paralleltest // the integration test have strict order
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			AddRoutes(r, &log, rest.ChiNamedURLParamsGetter, paymentService)

			srv := httptest.NewServer(r)
			defer srv.Close()

			req := httptest.NewRequest(tt.httpMethod, srv.URL+tt.urlPath, strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.expectedHTTPStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v, body: %v",
					status, http.StatusOK, rr.Body.String())
			}
		})
	}
}
