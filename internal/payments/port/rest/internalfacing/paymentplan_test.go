//nolint:cyclop // ignore for testing
package internalfacing

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"golangreferenceapi/internal/payments/service"

	"golangreferenceapi/internal/payments/port/rest"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func Test_createPendingPaymentPlanHandlerError(t *testing.T) {
	t.Parallel()

	type args struct {
		userUUID     string
		reqBody      io.Reader
		paramsGetter handlerwrap.NamedURLParamsGetter
	}

	tests := []struct {
		name              string
		args              args
		wantErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "returns 400 if passing a invalid param user_uuid",
			args: args{
				userUUID:     "x",
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: handlerwrap.ParsingParamError{
				Name:  urlParamUserUUID,
				Value: "x",
			}.ToErrorResponse(),
		},
		{
			name: "returns 400 if passing a broken reqBody",
			args: args{
				userUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				reqBody: readerFunc(func(p []byte) (int, error) {
					return 0, errors.New("failed")
				}),
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusBadRequest},
		},
		{
			name: "returns 400 if passing a invalid body",
			args: args{
				userUUID:     "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				reqBody:      strings.NewReader(`{x}`),
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusBadRequest},
		},
	}

	paymentService := service.NewPaymentPlanService()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", tt.args.reqBody)
			setURLParams(req, map[string]string{urlParamUserUUID: tt.args.userUUID})

			_, errRsp := createPendingPaymentPlanHandler(tt.args.paramsGetter, paymentService)(req)
			if errRsp != nil {
				if errRsp.HTTPStatusCode != tt.wantErrorResponse.HTTPStatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v", errRsp.HTTPStatusCode, tt.wantErrorResponse.HTTPStatusCode)
				}

				return
			}

			t.Errorf("expected error not exist")
		})
	}
}

func TestCreatePendingPaymentPlanHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		userUUID     string
		reqBody      io.Reader
		paramsGetter handlerwrap.NamedURLParamsGetter
	}

	arg := args{
		userUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
		reqBody: strings.NewReader(`{
					"payment": {
						"id": "03baa9e6-6ed6-4868-9ef9-b99c8452f270",
						"currency": "usdc",
						"total_amount": "100.0",
						"installments": [
							{ "id": "7b0547c6-2a32-47df-ab5d-0059ead32d2e", "due_at": "2022-06-01T14:02:03.000Z", "amount": "50", "currency": "usdc"},
							{ "id": "ca9407e7-bded-4caa-840d-c58573e3e6cf", "due_at": "2022-06-15T14:02:03.000Z", "amount": "50", "currency": "usdc"}
						]
				  }
				}`),
		paramsGetter: rest.ChiNamedURLParamsGetter,
	}

	wantResponse := &handlerwrap.Response{
		HTTPStatusCode: http.StatusOK,
		Body: CreatePendingPaymentPlanResponse{
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
		},
	}

	paymentService := service.NewPaymentPlanService()

	req := httptest.NewRequest("POST", "/", arg.reqBody)
	setURLParams(req, map[string]string{urlParamUserUUID: arg.userUUID})

	resp, errRsp := createPendingPaymentPlanHandler(arg.paramsGetter, paymentService)(req)
	if errRsp != nil {
		t.Errorf("returned unexpected error code: %v, err: %v", errRsp.HTTPStatusCode, errRsp.Error.Error())
	}

	pendingPlanExpect, ok := wantResponse.Body.(CreatePendingPaymentPlanResponse)
	if !ok {
		t.Errorf("failed to assertion")
	}

	pendingPlanActual, ok := resp.Body.(CreatePendingPaymentPlanResponse)
	if !ok {
		t.Errorf("failed to assertion")
	}

	switch {
	case pendingPlanActual.PendingPayment.ID != pendingPlanExpect.PendingPayment.ID:
		t.Errorf("expect: %v, actual: %v",
			pendingPlanExpect.PendingPayment.ID,
			pendingPlanActual.PendingPayment.ID,
		)
	case pendingPlanActual.PendingPayment.Currency != pendingPlanExpect.PendingPayment.Currency:
		t.Errorf("expect: %v, actual: %v",
			pendingPlanExpect.PendingPayment.Currency,
			pendingPlanActual.PendingPayment.Currency,
		)
	case pendingPlanActual.PendingPayment.Status != pendingPlanExpect.PendingPayment.Status:
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

		switch {
		case actualInstallment.Currency != expectedInstallment.Currency:
			t.Errorf("expect: %v, actual: %v",
				actualInstallment.Currency,
				expectedInstallment.Currency,
			)
		case actualInstallment.Amount != expectedInstallment.Amount:
			t.Errorf("expect: %v, actual: %v",
				actualInstallment.Amount,
				expectedInstallment.Amount,
			)
		}
	}
}

func Benchmark_createPendingPaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{urlParamUserUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5"})

	paymentService := service.NewPaymentPlanService()

	h := createPendingPaymentPlanHandler(rest.ChiNamedURLParamsGetter, paymentService)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}

func Test_cancelPaymentPlanHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		userUUID        string
		paymentPlanUUID string
		paramsGetter    handlerwrap.NamedURLParamsGetter
	}

	tests := []struct {
		name              string
		args              args
		wantResponse      *handlerwrap.Response
		wantErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "happy path",
			args: args{
				userUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paymentPlanUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paramsGetter:    rest.ChiNamedURLParamsGetter,
			},
			wantResponse: &handlerwrap.Response{
				HTTPStatusCode: http.StatusOK,
			},
		},
		{
			name: "returns 400 if passing a invalid param user_uuid",
			args: args{
				userUUID:     "x",
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: handlerwrap.ParsingParamError{
				Name:  urlParamUserUUID,
				Value: "x",
			}.ToErrorResponse(),
		},
		{
			name: "returns 400 if passing a invalid param uuid",
			args: args{
				userUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paymentPlanUUID: "x",
				paramsGetter:    rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: handlerwrap.ParsingParamError{
				Name:  urlParamUserUUID,
				Value: "x",
			}.ToErrorResponse(),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", nil)
			setURLParams(req, map[string]string{
				urlParamUserUUID:        tt.args.userUUID,
				urlParamPaymentPlanUUID: tt.args.paymentPlanUUID,
			})

			resp, errRsp := cancelPaymentPlanHandler(tt.args.paramsGetter)(req)
			if tt.wantErrorResponse != nil {
				if errRsp.HTTPStatusCode != tt.wantErrorResponse.HTTPStatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v",
						errRsp.HTTPStatusCode, tt.wantErrorResponse.HTTPStatusCode)
				}

				if errRsp.ErrorCode != tt.wantErrorResponse.ErrorCode {
					t.Errorf("returned unexpected error code: got %v want %v",
						errRsp.ErrorCode, tt.wantErrorResponse.ErrorCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response: got %v want %v", resp, tt.wantResponse)
			}
		})
	}
}

func Benchmark_cancelPaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{
		urlParamUserUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
		urlParamPaymentPlanUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
	})

	h := cancelPaymentPlanHandler(rest.ChiNamedURLParamsGetter)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}

func Test_completePaymentPlanHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		userUUID        string
		paymentPlanUUID string
		paramsGetter    handlerwrap.NamedURLParamsGetter
	}

	tests := []struct {
		name              string
		args              args
		wantResponse      *handlerwrap.Response
		wantErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "happy path",
			args: args{
				userUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paymentPlanUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paramsGetter:    rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusNotFound},
		},
		{
			name: "returns 400 if passing a invalid param user_uuid",
			args: args{
				userUUID:     "x",
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusBadRequest},
		},
		{
			name: "returns 400 if passing a invalid param uuid",
			args: args{
				userUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paymentPlanUUID: "x",
				paramsGetter:    rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusBadRequest},
		},
	}
	paymentService := service.NewPaymentPlanService()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", nil)
			setURLParams(req, map[string]string{
				urlParamUserUUID:        tt.args.userUUID,
				urlParamPaymentPlanUUID: tt.args.paymentPlanUUID,
				urlParamPaymentUUID:     tt.args.userUUID,
				urlParamInstallmentID:   tt.args.paymentPlanUUID,
			})

			resp, errRsp := completePaymentPlanHandler(tt.args.paramsGetter, paymentService)(req)
			if tt.wantErrorResponse != nil {
				if errRsp.HTTPStatusCode != tt.wantErrorResponse.HTTPStatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v",
						errRsp.HTTPStatusCode, tt.wantErrorResponse.HTTPStatusCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response: got %v want %v", resp, tt.wantResponse)
			}
		})
	}
}

func Benchmark_completePaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{
		urlParamUserUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
		urlParamPaymentPlanUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
	})

	paymentService := service.NewPaymentPlanService()
	h := completePaymentPlanHandler(rest.ChiNamedURLParamsGetter, paymentService)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}

func Test_refundHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		reqBody io.Reader
	}

	tests := []struct {
		name              string
		args              args
		wantResponse      *handlerwrap.Response
		wantErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "happy path",
			args: args{
				reqBody: strings.NewReader(`{}`),
			},
			wantResponse: &handlerwrap.Response{
				HTTPStatusCode: http.StatusOK,
			},
		},
		{
			name: "returns 400 if passing a broken reqBody",
			args: args{
				reqBody: readerFunc(func(p []byte) (int, error) {
					return 0, errors.New("failed")
				}),
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusBadRequest},
		},

		{
			name: "returns 400 if passing a invalid reqBody",
			args: args{
				reqBody: strings.NewReader(`{x}`),
			},
			wantErrorResponse: &handlerwrap.ErrorResponse{HTTPStatusCode: http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", tt.args.reqBody)

			resp, errRsp := refundHandler()(req)
			if tt.wantErrorResponse != nil {
				if errRsp.HTTPStatusCode != tt.wantErrorResponse.HTTPStatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v", errRsp.HTTPStatusCode, tt.wantErrorResponse.HTTPStatusCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response: got %v want %v", resp, tt.wantResponse)
			}
		})
	}
}

func Benchmark_refundHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{urlParamUserUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5"})

	h := refundHandler()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
