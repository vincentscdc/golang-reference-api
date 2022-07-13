package userfacing

import (
	"context"
	"reflect"
	"testing"

	"golangreferenceapi/internal/payments/transport/grpc/bnplapi/creditline/v1"
)

func TestPayLaterServer_GetCreditLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *creditline.GetCreditLineRequest
		wantResponse *creditline.GetCreditLineResponse
	}{
		{
			name:    "SimpleTest",
			request: &creditline.GetCreditLineRequest{},
			wantResponse: &creditline.GetCreditLineResponse{
				CreditInfo: &creditline.CreditInfo{
					TotalAmount:     "1000",
					AvailableAmount: "1000",
					Currency:        "USDC",
					Status:          "active",
				},
				Error: nil,
			},
		},
	}
	server := NewPayLaterServer()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp, err := server.GetCreditLine(context.Background(), tt.request)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			} else if !reflect.DeepEqual(resp, tt.wantResponse) {
				t.Errorf("returned unexpected response %v", tt.wantResponse)
			}
		})
	}
}

func BenchmarkPayLaterServer_GetCreditLine(b *testing.B) {
	server := NewPayLaterServer()
	ctx := context.Background()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_, _ = server.GetCreditLine(ctx, &creditline.GetCreditLineRequest{})
	}
}
