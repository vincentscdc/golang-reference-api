package userfacing

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func Test_getPaymentPlansHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		query                 string
		expectedResponse      *handlerwrap.Response
		expectedErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name:  "happy path",
			query: "offset=10&limit=10&created_at_order=desc",
			expectedResponse: &handlerwrap.Response{
				Body:           ([]PaymentPlanResponse)(nil),
				HTTPStatusCode: http.StatusOK,
			},
		},
		{
			name:                  "invalid pagination query params",
			query:                 "offset=x&limit=x&created_at_order=x",
			expectedErrorResponse: handlerwrap.ParsingParamError{}.ToErrorResponse(),
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				req := httptest.NewRequest("GET", "/?"+tt.query, nil)
				setRequestHeaderUserID(req, "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5")

				resp, err := listPaymentPlansHandler()(req)
				if tt.expectedErrorResponse != nil {
					if err.ErrorCode != tt.expectedErrorResponse.ErrorCode {
						t.Errorf("returned a unexpected error code got %v want %v", err.ErrorCode, tt.expectedErrorResponse.ErrorCode)
					}

					return
				}
				if !reflect.DeepEqual(resp, tt.expectedResponse) {
					t.Errorf("returned a unexpected response got %#v want %#v", resp, tt.expectedResponse)
				}
			})
		})
	}
}

func Benchmark_getPaymentPlansHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	h := listPaymentPlansHandler()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
