package internalfacing

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bnpl/internal/port/rest"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func Test_getCreditLineHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		uuidStr      string
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
				uuidStr:      "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantResponse: &handlerwrap.Response{
				Body: &CreditLineResponse{
					Limit:    "1000",
					Balance:  "1000",
					Currency: "USDC",
					Status:   "active",
				},
				HTTPStatusCode: http.StatusOK,
			},
		},
		{
			name: "returns 400 when param user_uuid is empty",
			args: args{
				paramsGetter: rest.ChiNamedURLParamsGetter,
			},
			wantErrorResponse: handlerwrap.MissingParamError{
				Name: urlParamUserUUID,
			}.ToErrorResponse(),
		},
		{
			name: "returns 400 when param user_uuid is invalid",
			args: args{
				uuidStr:      "x",
				paramsGetter: rest.ChiNamedURLParamsGetter,
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

			req := httptest.NewRequest("GET", "/", nil)
			setURLParams(req, map[string]string{urlParamUserUUID: tt.args.uuidStr})

			resp, err := getCreditLineHandler(tt.args.paramsGetter)(req)
			if tt.wantErrorResponse != nil {
				if err.HTTPStatusCode != tt.wantErrorResponse.HTTPStatusCode {
					t.Errorf("returned unexpected HTTP status code: got %v want %v", err.HTTPStatusCode, tt.wantErrorResponse.HTTPStatusCode)
				}

				if err.ErrorCode != tt.wantErrorResponse.ErrorCode {
					t.Errorf("returned unexpected error code: got %v want %v", err.ErrorCode, tt.wantErrorResponse.ErrorCode)
				}

				return
			}

			if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response: got %v want %v", resp, tt.wantResponse)
			}
		})
	}
}

func Benchmark_getCreditLineHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	setURLParams(req, map[string]string{urlParamUserUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5"})

	h := getCreditLineHandler(rest.ChiNamedURLParamsGetter)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
