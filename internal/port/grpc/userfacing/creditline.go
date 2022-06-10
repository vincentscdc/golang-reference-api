package userfacing

import (
	"context"

	"bnpl/internal/port/grpc/protos/creditline/v1"
)

var _ creditline.PayLaterServiceServer = (*PayLaterServer)(nil)

type PayLaterServer struct {
	creditline.UnimplementedPayLaterServiceServer
}

func NewPayLaterServer() *PayLaterServer {
	return &PayLaterServer{}
}

func (s *PayLaterServer) GetCreditLine(ctx context.Context, request *creditline.GetCreditLineRequest) (
	*creditline.GetCreditLineResponse, error,
) {
	resp := &creditline.GetCreditLineResponse{
		CreditInfo: &creditline.CreditInfo{
			TotalAmount:     "1000",
			AvailableAmount: "1000",
			Currency:        "USDC",
			Status:          "active",
		},
		Error: nil,
	}

	return resp, nil
}
