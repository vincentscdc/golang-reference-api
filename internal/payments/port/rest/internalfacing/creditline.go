package internalfacing

import (
	"net/http"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

type CreditLineResponse struct {
	Limit    string `json:"limit"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

// getCreditLineHandler renders the credit line
// @Summary Renders a user's credit line info
// @Description returns the amount and status
// @Tags credit_line
// @Produce json
// @Router /api/internal/pay_later/user/{user_uuid}/credit_line [get]
// @Param user_uuid path string true "User UUID"
// @Success 200 {object} CreditLineResponse
func getCreditLineHandler(paramsGetter handlerwrap.NamedURLParamsGetter) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		// get params
		_, err := parseUUIDFormatParam(req.Context(), paramsGetter, urlParamUserUUID)
		if err != nil {
			return nil, err
		}

		info := &CreditLineResponse{
			Limit:    "1000",
			Balance:  "1000",
			Currency: "USDC",
			Status:   "active",
		}
		resp := &handlerwrap.Response{Body: info, HTTPStatusCode: http.StatusOK}

		return resp, nil
	}
}
