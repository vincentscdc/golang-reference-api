package userfacing

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"golangreferenceapi/internal/payments/common"
	"golangreferenceapi/internal/payments/mock/servicemock"
	"golangreferenceapi/internal/payments/service"
	"golangreferenceapi/internal/payments/transport/rest"

	"github.com/monacohq/golang-common/transport/http/middleware/cryptouseruuid"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
)

func Test_listPaymentPlansHandler_SuccessCase(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV4())
	installmentID := uuid.Must(uuid.NewV4())
	expectedResult := []service.PaymentPlans{{
		ID:          userID.String(),
		Currency:    "usdc",
		TotalAmount: "100",
		Status:      "pending",
		CreatedAt:   time.Now().Format(common.TimeFormat),
		Installments: []service.PaymentPlanInstallment{{
			ID:       installmentID.String(),
			Amount:   "100",
			Currency: "usdc",
			DueAt:    time.Now().Format(common.TimeFormat),
			Status:   "pending",
		}},
	}}
	resultBody := PaymentPlanResponse{Payments: expectedResult}
	tests := []struct {
		name                  string
		query                 string
		expectedResponse      *handlerwrap.Response
		expectedErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name:  "happy path",
			query: "offset=0&limit=10&created_at_order=desc",
			expectedResponse: &handlerwrap.Response{
				Body:       resultBody,
				StatusCode: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			paymentService := servicemock.NewMockPaymentPlanService(mockCtrl)

			gomock.InOrder(
				paymentService.EXPECT().GetPaymentPlanByUserID(gomock.Any(), userID).
					Return(expectedResult, nil).AnyTimes(),
			)

			req := httptest.NewRequest("GET", "/?"+tt.query, nil)
			req = req.WithContext(cryptouseruuid.SetUserUUID(req.Context(), &userID))

			resp, err := listPaymentPlansHandler(paymentService)(req)
			if tt.expectedErrorResponse != nil {
				if err.StatusCode != tt.expectedErrorResponse.StatusCode {
					t.Errorf("returned a unexpected error code got %v want %v", err.StatusCode, tt.expectedErrorResponse.StatusCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.expectedResponse) {
				t.Errorf("returned a unexpected response got %#v want %#v", resp, tt.expectedResponse)
			}
		})
	}
}

func Test_listPaymentPlansHandler_ParamsError(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name                  string
		query                 string
		expectedResponse      *handlerwrap.Response
		expectedErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name:                  "invalid pagination query params",
			query:                 "offset=x&limit=x&created_at_order=x",
			expectedErrorResponse: handlerwrap.ParsingParamError{}.ToErrorResponse(),
		},
	}

	mockCtrl := gomock.NewController(t)
	paymentService := servicemock.NewMockPaymentPlanService(mockCtrl)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "/?"+tt.query, nil)

			req = req.WithContext(cryptouseruuid.SetUserUUID(req.Context(), &userID))

			resp, err := listPaymentPlansHandler(paymentService)(req)
			if tt.expectedErrorResponse != nil {
				if err.StatusCode != tt.expectedErrorResponse.StatusCode {
					t.Errorf("returned a unexpected error code got %v want %v", err.StatusCode, tt.expectedErrorResponse.StatusCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.expectedResponse) {
				t.Errorf("returned a unexpected response got %#v want %#v", resp, tt.expectedResponse)
			}
		})
	}
}

func Test_listPaymentPlansHandler_InternalError(t *testing.T) {
	t.Parallel()

	userID := uuid.Must(uuid.NewV4())

	tests := []struct {
		uid              uuid.UUID
		name             string
		expectedResponse *handlerwrap.Response
		err              error
	}{
		{
			uid:  uuid.Must(uuid.NewV4()),
			name: "function internal error",
			err:  service.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				paymentService := servicemock.NewMockPaymentPlanService(mockCtrl)

				gomock.InOrder(
					paymentService.EXPECT().
						GetPaymentPlanByUserID(gomock.Any(), userID).
						Return(nil, tt.err).AnyTimes(),
				)

				req := httptest.NewRequest("GET", "/", nil)

				req = req.WithContext(cryptouseruuid.SetUserUUID(req.Context(), &userID))

				resp, err := listPaymentPlansHandler(paymentService)(req)
				if resp != nil {
					t.Errorf("unexpected response %v", resp)
				}

				if !reflect.DeepEqual(rest.ServiceErrorToErrorResp(tt.err), err) { // nolint: deepequalerrors // linter bug these are responses, not errors
					t.Errorf("unexpected error, expected: %v, actual: %v", rest.ServiceErrorToErrorResp(tt.err), err)
				}
			})
		})
	}
}

func Test_listPaymentPlansHandler_MissingUserID(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	paymentService := servicemock.NewMockPaymentPlanService(mockCtrl)

	req := httptest.NewRequest("GET", "/", nil)

	_, err := listPaymentPlansHandler(paymentService)(req)
	if err != nil && !errors.As(err.Error, &cryptouseruuid.UserIDNotFoundError{}) {
		t.Errorf("unexpected status code expect: %v, actual: %v", err.ErrorCode, http.StatusBadRequest)
	}
}

func Benchmark_listPaymentPlansHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	paymentService := service.NewPaymentPlanService()
	h := listPaymentPlansHandler(paymentService)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
