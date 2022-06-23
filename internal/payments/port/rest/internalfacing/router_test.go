package internalfacing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"golangreferenceapi/internal/payments/port/rest"
	"golangreferenceapi/internal/payments/service"
)

func TestAddRoutes(t *testing.T) {
	t.Parallel()

	log := zerolog.Nop().With().Logger()

	tests := []struct {
		name                   string
		httpMethod             string
		urlPath                string
		reqBody                string
		expectedHTTPStatusCode int
	}{
		{
			name:                   "happy path for getting credit line",
			httpMethod:             "GET",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/credit_line",
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:       "happy path for creating a pending payment plan",
			httpMethod: "POST",
			urlPath:    "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/payment_plans",
			reqBody: `{
					"payment": {
						"id": "03baa9e6-6ed6-4868-9ef9-b99c8452f270",
						"currency": "usdc",
						"total_amount": "100.0",
						"installments": [
							{ "id": "7b0547c6-2a32-47df-ab5d-0059ead32d2e", "due_at": "2022-06-01T14:02:03.000Z", "amount": "50", "currency": "usdc"},
							{ "id": "ca9407e7-bded-4caa-840d-c58573e3e6cf", "due_at": "2022-06-15T14:02:03.000Z", "amount": "50", "currency": "usdc"}
						]
				  }
				}`,
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "happy path for canceling a pending payment plan",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/payment_plans/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/cancel",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "happy path for completing payment plan",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/payment_plans/03baa9e6-6ed6-4868-9ef9-b99c8452f270/installments/7b0547c6-2a32-47df-ab5d-0059ead32d2e/payment",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			name:                   "happy path for refunding payment",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/refund",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
		},
	}

	paymentService := service.NewPaymentPlanService()

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
