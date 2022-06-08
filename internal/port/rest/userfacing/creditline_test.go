package userfacing

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func Test_getCreditLineHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		expectedResponse      *handlerwrap.Response
		expectedErrorResponse *handlerwrap.ErrorResponse
	}{
		{
			name: "happy path",
			expectedResponse: &handlerwrap.Response{
				Body: &CreditLineResponse{
					TotalAmount:     "1000",
					AvailableAmount: "1000",
					Currency:        "USDC",
					Status:          "active",
				},
				HTTPStatusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", "/", nil)
			setRequestHeaderUserID(req, "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5")

			resp, err := getCreditLineHandler()(req)
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
	}
}

func Benchmark_getCreditLineHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)

	h := getCreditLineHandler()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = h(req)
	}
}
