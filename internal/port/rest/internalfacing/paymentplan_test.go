package internalfacing

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"golangreferenceapi/internal/port/rest"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func Test_createPendingPaymentPlanHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		userUUID     string
		reqBody      io.Reader
		paramsGetter handlerwrap.NamedURLParamsGetter
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
				userUUID:     "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				reqBody:      strings.NewReader(`{}`),
				paramsGetter: rest.ChiNamedURLParamsGetter,
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
			name: "returns 400 if passing a broken reqBody",
			args: args{
				userUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				reqBody: readerFunc(func(p []byte) (int, error) {
					return 0, errors.New("failed")
				}),
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: InvalidRequestBodyError{}.ToErrorResponse(),
		},
		{
			name: "returns 400 if passing a invalid body",
			args: args{
				userUUID:     "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				reqBody:      strings.NewReader(`{x}`),
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: InvalidRequestBodyError{}.ToErrorResponse(),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/", tt.args.reqBody)
			setURLParams(req, map[string]string{urlParamUserUUID: tt.args.userUUID})

			resp, errRsp := createPendingPaymentPlanHandler(tt.args.paramsGetter)(req)
			if tt.wantErrorResponse != nil {
				if errRsp.HTTPStatusCode != tt.wantErrorResponse.HTTPStatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v", errRsp.HTTPStatusCode, tt.wantErrorResponse.HTTPStatusCode)
				}

				if errRsp.ErrorCode != tt.wantErrorResponse.ErrorCode {
					t.Errorf("returned unexpected error code: got %v want %v", errRsp.ErrorCode, tt.wantErrorResponse.ErrorCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response: got %v want %v", resp, tt.wantResponse)
			}
		})
	}
}

func Benchmark_createPendingPaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{urlParamUserUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5"})

	h := createPendingPaymentPlanHandler(rest.ChiNamedURLParamsGetter)

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

			resp, errRsp := completePaymentPlanHandler(tt.args.paramsGetter)(req)
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

func Benchmark_completePaymentPlanHandler(b *testing.B) {
	req := httptest.NewRequest("POST", "/", nil)
	setURLParams(req, map[string]string{
		urlParamUserUUID:        "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
		urlParamPaymentPlanUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
	})

	h := completePaymentPlanHandler(rest.ChiNamedURLParamsGetter)

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
			wantErrorResponse: InvalidRequestBodyError{}.ToErrorResponse(),
		},

		{
			name: "returns 400 if passing a invalid reqBody",
			args: args{
				reqBody: strings.NewReader(`{x}`),
			},
			wantErrorResponse: InvalidRequestBodyError{}.ToErrorResponse(),
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

				if errRsp.ErrorCode != tt.wantErrorResponse.ErrorCode {
					t.Errorf("returned unexpected error code: got %v want %v", errRsp.ErrorCode, tt.wantErrorResponse.ErrorCode)
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
