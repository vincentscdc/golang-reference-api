package userfacing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golangreferenceapi/internal/payments/mock/servicemock"
	"golangreferenceapi/internal/payments/port/rest"
	"golangreferenceapi/internal/payments/service"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
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
	}{
		{
			name:                   "happy path for payment_plans",
			httpMethod:             "GET",
			urlPath:                "/api/pay_later/payment_plans",
			expectedHTTPStatusCode: http.StatusOK,
		},
	}

	userID := uuid.Must(uuid.NewV4())

	paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))
	paymentService.EXPECT().GetPaymentPlanByUserID(gomock.Any(), userID).Return([]service.PaymentPlans{}, nil)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			AddRoutes(r, &log, paymentService)

			srv := httptest.NewServer(r)
			defer srv.Close()

			req := httptest.NewRequest(tt.httpMethod, srv.URL+tt.urlPath, nil)
			setRequestHeaderUserID(req, userID.String())
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tt.expectedHTTPStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		})
	}
}

func setRequestHeaderUserID(r *http.Request, uuid string) {
	r.Header.Set(rest.HTTPHeaderKeyUserUUID, uuid)
}
