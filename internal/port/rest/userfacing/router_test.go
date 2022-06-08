package userfacing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		expectedHTTPStatusCode int
		expectedBody           string
	}{
		{
			name:                   "happy path for credit line",
			httpMethod:             "GET",
			urlPath:                "/api/pay_later/credit_line",
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `{"credit_info":{"total_amount":"1000","available_amount":"1000","currency":"USDC","status":"active"},"ok":true}`,
		},
		{
			name:                   "happy path for payment plans",
			httpMethod:             "GET",
			urlPath:                "/api/pay_later/payment_plans",
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           `{"ok":true,"payment_plans":null}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			AddRoutes(r, &log)

			srv := httptest.NewServer(r)
			defer srv.Close()

			req := httptest.NewRequest(tt.httpMethod, srv.URL+tt.urlPath, nil)
			setRequestHeaderUserID(req, "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5")
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
