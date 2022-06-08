package internalfacing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bnpl/internal/port/rest"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
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
		expectedBody           string
	}{
		{
			name:                   "happy path for getting credit line",
			httpMethod:             "GET",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/credit_line",
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `{"limit":"1000","balance":"1000","currency":"USDC","status":"active"}`,
		},
		{
			name:                   "happy path for creating a pending payment plan",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/payment_plans",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `null`,
		},
		{
			name:                   "happy path for canceling a pending payment plan",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/payment_plans/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/cancel",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `null`,
		},
		{
			name:                   "happy path for completing payment plan",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/user/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/payment_plans/b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5/complete",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `null`,
		},
		{
			name:                   "happy path for refunding payment",
			httpMethod:             "POST",
			urlPath:                "/api/internal/pay_later/refund",
			reqBody:                `{}`,
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `null`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			AddRoutes(r, &log, rest.ChiNamedURLParamsGetter)

			srv := httptest.NewServer(r)
			defer srv.Close()

			req := httptest.NewRequest(tt.httpMethod, srv.URL+tt.urlPath, strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
