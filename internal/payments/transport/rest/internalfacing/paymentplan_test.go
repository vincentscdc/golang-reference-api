package internalfacing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"

	"golangreferenceapi/internal/payments/common"
	"golangreferenceapi/internal/payments/mock/servicemock"
	"golangreferenceapi/internal/payments/service"
	"golangreferenceapi/internal/payments/transport/rest"
)

func Test_createPendingPaymentPlanHandlerInputError(t *testing.T) {
	t.Parallel()

	type args struct {
		reqBody io.Reader
	}

	tests := []struct {
		name              string
		args              args
		wantErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "returns 400 if passing a broken reqBody",
			args: args{
				reqBody: readerFunc(func(p []byte) (int, error) {
					return 0, errors.New("failed")
				}),
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{StatusCode: http.StatusBadRequest},
		},
		{
			name: "returns 400 if passing a invalid body",
			args: args{
				reqBody: strings.NewReader(`{x}`),
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{StatusCode: http.StatusBadRequest},
		},
	}

	paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", tt.args.reqBody)

			_, errRsp := createPendingPaymentPlanHandler(paymentService)(req)
			if errRsp != nil {
				if errRsp.StatusCode != tt.wantErrorResponse.StatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v", errRsp.StatusCode, tt.wantErrorResponse.StatusCode)
				}

				return
			}

			t.Errorf("expected error not exist")
		})
	}
}

func Test_completePaymentPlanHandlerUserIDNotFound(t *testing.T) {
	t.Parallel()

	paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))

	tests := []struct {
		name   string
		hander handlerwrap.TypedHandler
	}{
		{
			name:   "complete payment handler",
			hander: completePaymentPlanHandler(rest.ChiNamedURLParamsGetter, paymentService),
		},
		// {
		// 	name:   "create pending payment handler",
		// 	hander: createPendingPaymentPlanHandler(paymentService),
		// },
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))

			resp, respErr := tt.hander(req)
			if resp != nil {
				t.Errorf("unexpected resp, %v", resp)
			}

			if reflect.TypeOf(respErr).Name() != reflect.TypeOf(handlerwrap.ParsingParamError{}.ToErrorResponse()).Name() {
				t.Errorf("unexpected error response, expected: %v, actual: %v",
					reflect.TypeOf(handlerwrap.ParsingParamError{}.ToErrorResponse()).Name(),
					reflect.TypeOf(respErr).Name(),
				)
			}
		})
	}
}

func Test_createPendingPaymentPlanHandler(t *testing.T) {
	t.Parallel()

	var (
		userUUID = uuid.Must(uuid.NewV4())
		planUUID = uuid.Must(uuid.NewV4())
		request  = CreatePendingPaymentPlanRequest{
			PendingPayment: service.CreatePaymentPlanParams{
				ID:          planUUID,
				UserID:      userUUID,
				Currency:    "usdc",
				TotalAmount: "100",
				Installments: []service.PaymentPlanInstallmentParams{
					{
						Amount:   "100",
						Currency: "usdc",
						DueAt:    time.Now().UTC().Round(time.Second),
					},
				},
			},
		}
		response = CreatePendingPaymentPlanResponse{
			PendingPayment: service.PaymentPlans{
				ID:          "03baa9e6-6ed6-4868-9ef9-b99c8452f270",
				Currency:    "usdc",
				TotalAmount: "200.0",
				Status:      service.PaymentInstallmentStatusPending,
				Installments: []service.PaymentPlanInstallment{
					{
						DueAt:    "2022-06-01T14:02:03.000Z",
						Amount:   "50",
						Currency: "usdc",
					},
					{
						DueAt:    "2022-06-01T14:02:03.000Z",
						Amount:   "50",
						Currency: "usdc",
					},
				},
			},
		}
		wantResponse = &handlerwrap.Response{
			StatusCode: http.StatusOK,
			Body:       response,
		}
	)

	paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))
	gomock.InOrder(
		paymentService.EXPECT().CreatePendingPaymentPlan(
			gomock.Any(),
			gomock.Eq(&request.PendingPayment),
		).Return(&response.PendingPayment, nil),
	)

	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Errorf("failed to unmarshal json")
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

	resp, errRsp := createPendingPaymentPlanHandler(paymentService)(req)
	if errRsp != nil {
		t.Errorf("returned unexpected error code: %v, err: %v", errRsp.StatusCode, errRsp.Error.Error())
	}

	pendingPlanExpect, ok := wantResponse.Body.(CreatePendingPaymentPlanResponse)
	if !ok {
		t.Errorf("failed to assertion")
	}

	pendingPlanActual, ok := resp.Body.(CreatePendingPaymentPlanResponse)
	if !ok {
		t.Errorf("failed to assertion")
	}

	if pendingPlanActual.PendingPayment.ID != pendingPlanExpect.PendingPayment.ID {
		t.Errorf("expect: %v, actual: %v",
			pendingPlanExpect.PendingPayment.ID,
			pendingPlanActual.PendingPayment.ID,
		)
	}

	if pendingPlanActual.PendingPayment.Currency != pendingPlanExpect.PendingPayment.Currency {
		t.Errorf("expect: %v, actual: %v",
			pendingPlanExpect.PendingPayment.Currency,
			pendingPlanActual.PendingPayment.Currency,
		)
	}

	if pendingPlanActual.PendingPayment.Status != pendingPlanExpect.PendingPayment.Status {
		t.Errorf("expect: %v, actual: %v",
			pendingPlanExpect.PendingPayment.Status,
			pendingPlanActual.PendingPayment.Status,
		)
	}

	if len(pendingPlanActual.PendingPayment.Installments) !=
		len(pendingPlanExpect.PendingPayment.Installments) {
		t.Errorf(
			"expected: %v, actual: %v",
			len(pendingPlanExpect.PendingPayment.Installments),
			len(pendingPlanActual.PendingPayment.Installments),
		)
	}

	for idx := range pendingPlanActual.PendingPayment.Installments {
		actualInstallment := pendingPlanActual.PendingPayment.Installments[idx]
		expectedInstallment := pendingPlanExpect.PendingPayment.Installments[idx]

		if actualInstallment.Currency != expectedInstallment.Currency {
			t.Errorf("expect: %v, actual: %v",
				actualInstallment.Currency,
				expectedInstallment.Currency,
			)
		}

		if actualInstallment.Amount != expectedInstallment.Amount {
			t.Errorf("expect: %v, actual: %v",
				actualInstallment.Amount,
				expectedInstallment.Amount,
			)
		}
	}
}

func Test_createPendingPaymentPlanHandlerServiceError(t *testing.T) {
	t.Parallel()

	var (
		userUUID = uuid.Must(uuid.NewV4())
		planUUID = uuid.Must(uuid.NewV4())
		request  = CreatePendingPaymentPlanRequest{
			PendingPayment: service.CreatePaymentPlanParams{
				ID:          planUUID,
				UserID:      userUUID,
				Currency:    "usdc",
				TotalAmount: "100",
				Installments: []service.PaymentPlanInstallmentParams{
					{
						Amount:   "100",
						Currency: "usdc",
						DueAt:    time.Now().UTC().Round(time.Second),
					},
				},
			},
		}
	)

	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "record not found",
			args: args{
				err: service.ErrRecordNotFound,
			},
		},
		{
			name: "failed to generate uuid",
			args: args{
				err: service.ErrGenerateUUID,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))
			gomock.InOrder(
				paymentService.EXPECT().CreatePendingPaymentPlan(
					gomock.Any(),
					gomock.Eq(&request.PendingPayment),
				).Return(nil, tt.args.err),
			)

			wantResponse := rest.ServiceErrorToErrorResp(tt.args.err)

			reqBody, err := json.Marshal(request)
			if err != nil {
				t.Errorf("failed to unmarshal json")
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

			resp, errRsp := createPendingPaymentPlanHandler(paymentService)(req)
			if resp != nil {
				t.Errorf("returned unexpected response: %v", resp)
			}

			if !reflect.DeepEqual(wantResponse, errRsp) { // nolint: deepequalerrors // linter bug these are responses, not errors
				t.Errorf("returned unexpected err. expected: %v, actual: %v", wantResponse, errRsp)
			}
		})
	}
}

func Benchmark_createPendingPaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)

	paymentService := service.NewPaymentPlanService()

	h := createPendingPaymentPlanHandler(paymentService)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}

func Test_completePaymentPlanHandler(t *testing.T) {
	t.Parallel()

	var (
		userUUID      = uuid.Must(uuid.NewV4())
		paymentPlanID = uuid.Must(uuid.NewV4())
		paramsGetter  = rest.ChiNamedURLParamsGetter
		payment       = service.PaymentPlans{
			ID:           paymentPlanID.String(),
			Currency:     "usdc",
			TotalAmount:  "100",
			Status:       "pending",
			CreatedAt:    time.Now().Format(common.TimeFormat),
			Installments: nil,
		}
		wantResponse = &handlerwrap.Response{
			StatusCode: http.StatusOK,
			Body:       CompletePaymentPlanResponse{Payment: payment},
		}
		request = CompletePaymentPlanRequest{
			Payment: service.CompletePaymentPlanParams{
				UserID: userUUID,
			},
		}
	)

	paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))

	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Errorf("failed to unmarshal json")
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

	setURLParams(req, map[string]string{
		urlParamPaymentUUID: paymentPlanID.String(),
	})

	gomock.InOrder(
		paymentService.EXPECT().CompletePaymentPlanCreation(
			gomock.Eq(req.Context()),
			gomock.Eq(paymentPlanID),
			gomock.Eq(&request.Payment),
		).Return(&payment, nil),
	)

	resp, errRsp := completePaymentPlanHandler(paramsGetter, paymentService)(req)
	if errRsp != nil {
		t.Errorf("returned unexpected error response: %v", errRsp)
	}

	if !reflect.DeepEqual(resp, wantResponse) {
		t.Errorf("returned unexpected err. expected: %v, actual: %v", wantResponse, resp)
	}
}

func Test_completePaymentPlanHandlerParamsError(t *testing.T) {
	t.Parallel()

	userUUID := uuid.Must(uuid.NewV4())

	type args struct {
		paymentPlanUUID string
		paramsGetter    handlerwrap.NamedURLParamsGetter
		request         CompletePaymentPlanRequest
	}

	tests := []struct {
		name              string
		args              args
		wantResponse      *handlerwrap.Response
		wantErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "returns 400 if passing a invalid param uuid",
			args: args{
				paymentPlanUUID: "x",
				paramsGetter:    rest.ChiNamedURLParamsGetter,
				request: CompletePaymentPlanRequest{
					Payment: service.CompletePaymentPlanParams{
						UserID: userUUID,
					},
				},
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{StatusCode: http.StatusBadRequest},
		},
		{
			name: "returns 400 if passing a invalid param uuid",
			args: args{
				paymentPlanUUID: "x",
				paramsGetter: func(ctx context.Context, key string) (string, *handlerwrap.ErrorResponse) {
					return "", handlerwrap.MissingParamError{Name: key}.ToErrorResponse()
				},
				request: CompletePaymentPlanRequest{
					Payment: service.CompletePaymentPlanParams{
						UserID: userUUID,
					},
				},
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{StatusCode: http.StatusBadRequest},
		},
	}

	ctrl := gomock.NewController(t)
	paymentService := servicemock.NewMockPaymentPlanService(ctrl)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.args.request)
			if err != nil {
				t.Errorf("failed to unmarshal json")
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

			setURLParams(req, map[string]string{
				urlParamPaymentUUID: tt.args.paymentPlanUUID,
			})

			resp, errRsp := completePaymentPlanHandler(tt.args.paramsGetter, paymentService)(req)
			if tt.wantErrorResponse != nil {
				if errRsp.StatusCode != tt.wantErrorResponse.StatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v",
						errRsp.StatusCode, tt.wantErrorResponse.StatusCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response: got %v want %v", resp, tt.wantResponse)
			}
		})
	}
}

func Test_completePaymentPlanHandlerServiceError(t *testing.T) {
	t.Parallel()

	var (
		userUUID      = uuid.Must(uuid.NewV4())
		paymentPlanID = uuid.Must(uuid.NewV4())
		paramsGetter  = rest.ChiNamedURLParamsGetter
		request       = CompletePaymentPlanRequest{
			Payment: service.CompletePaymentPlanParams{
				UserID: userUUID,
			},
		}
	)

	type args struct {
		err error
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "record not found",
			args: args{err: service.ErrRecordNotFound},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wantResponse := rest.ServiceErrorToErrorResp(tt.args.err)

			paymentService := servicemock.NewMockPaymentPlanService(gomock.NewController(t))
			gomock.InOrder(
				paymentService.EXPECT().CompletePaymentPlanCreation(
					gomock.Any(),
					gomock.Eq(paymentPlanID),
					gomock.Eq(&request.Payment),
				).Return(nil, service.ErrRecordNotFound),
			)

			reqBody, err := json.Marshal(request)
			if err != nil {
				t.Errorf("failed to unmarshal json")
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

			setURLParams(req, map[string]string{
				urlParamPaymentUUID: paymentPlanID.String(),
			})

			resp, errRsp := completePaymentPlanHandler(paramsGetter, paymentService)(req)
			if resp != nil {
				t.Errorf("returned unexpected response: %v", resp)
			}

			if !reflect.DeepEqual(wantResponse, errRsp) { // nolint: deepequalerrors // linter bug these are responses, not errors
				t.Errorf("returned unexpected err. expected: %v, actual: %v", wantResponse, errRsp)
			}
		})
	}
}

func Benchmark_completePaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{
		urlParamPaymentUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
	})

	paymentService := service.NewPaymentPlanService()
	h := completePaymentPlanHandler(rest.ChiNamedURLParamsGetter, paymentService)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
