package userfacing

import (
	"net/http"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

type CreditLineResponse struct {
	TotalAmount     string `json:"total_amount"`
	AvailableAmount string `json:"available_amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
}

func getCreditLineHandler() handlerwrap.TypedHandler {
	return func(r *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		// user uuid
		_ = getUserUUID(r.Context())

		info := &CreditLineResponse{
			TotalAmount:     "1000",
			AvailableAmount: "1000",
			Currency:        "USDC",
			Status:          "active",
		}
		resp := &handlerwrap.Response{Body: info, HTTPStatusCode: http.StatusOK}

		return resp, nil
	}
}
