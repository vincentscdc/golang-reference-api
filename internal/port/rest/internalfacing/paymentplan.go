package internalfacing

import (
	"io"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

type PendingPayment struct {
	ID     string `json:"id"`
	Output string `json:"output"`
	Value  string `json:"value"`
	Meta   struct {
		Currency          string `json:"currency"`
		Amount            string `json:"amount"`
		Recipient         string `json:"recipient"`
		Items             string `json:"items"`
		CustomID          string `json:"custom_id"`
		QuotationID       string `json:"quotation_id"`
		ValueFiat         string `json:"value_fiat"`
		MerchantReference string `json:"merchant_reference"`
		LiveMode          bool   `json:"live_mode"`
		Status            string `json:"status"`
		RemainingTime     string `json:"remaining_time"`
		Deadline          string `json:"deadline"`
		CryptoCurrency    string `json:"crypto_currency"`
		IsApproved        bool   `json:"is_approved"`
		CryptoAmounts     struct {
			USDC string `json:"USDC"`
			CRO  string `json:"CRO"`
		} `json:"crypto_amounts"`
		PayLaterInstallments []struct {
			Date     string `json:"date"`
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"pay_later_installments"`
		PayLaterAmount struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"pay_later_amount"`
	} `json:"meta"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

type CreatePendingPaymentPlanRequest struct {
	PendingPayment     PendingPayment `json:"payment"`
	UserWalletCurrency string         `json:"user_wallet_currency"`
}

// createPendingPaymentPlanHandler creates a pending payment plan
// @Summary Creates a pending a payment plan
// @Description pre creates a payment plan
// @Tags payment_plan
// @Produce json
// @Router /api/internal/pay_later/user/{user_uuid}/payment_plans [post]
// @Param pre_create_payment_plan_request body CreatePendingPaymentPlanRequest true "Pre create payment plan reqBody"
// @Param user_uuid path string true "User UUID"
// @Success 200
// @Failure 400 {object} handlerwrap.ErrorResponse "bad reqBody"
// @Failure 500 {object} handlerwrap.ErrorResponse "internal error"
func createPendingPaymentPlanHandler(paramsGetter handlerwrap.NamedURLParamsGetter) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		var (
			userUUID *uuid.UUID
			request  CreatePendingPaymentPlanRequest
			err      *handlerwrap.ErrorResponse
		)

		userUUID, err = parseUUIDFormatParam(req.Context(), paramsGetter, urlParamUserUUID)
		if err != nil {
			return nil, err
		}

		_ = userUUID

		body, readErr := io.ReadAll(req.Body)
		if readErr != nil {
			return nil, InvalidRequestBodyError{Err: readErr}.ToErrorResponse()
		}

		if err := json.Unmarshal(body, &request); err != nil {
			return nil, InvalidRequestBodyError{Err: err, Data: string(body)}.ToErrorResponse()
		}

		return &handlerwrap.Response{
			Body:           nil,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}

// cancelPaymentPlanHandler cancels a payment plan
// @Summary Cancels a payment plan
// @Description cancels a payment plan
// @Tags payment_plan
// @Produce json
// @Router /api/internal/pay_later/user/{user_uuid}/payment_plans/{uuid}/cancel [post]
// @Param user_uuid path string true "User UUID"
// @Param uuid path string true "Payment Plan UUID"
// @Success 200
// @Failure 400 {object} handlerwrap.ErrorResponse "bad reqBody"
// @Failure 401 {object} handlerwrap.ErrorResponse "payment plan not belongs to user"
// @Failure 403 {object} handlerwrap.ErrorResponse "payment plan is not in pending"
// @Failure 404 {object} handlerwrap.ErrorResponse "payment plan not found"
// @Failure 500 {object} handlerwrap.ErrorResponse "internal error"
func cancelPaymentPlanHandler(paramsGetter handlerwrap.NamedURLParamsGetter) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		_, err := parsePaymentPlanParam(req.Context(), paramsGetter)
		if err != nil {
			return nil, err
		}

		return &handlerwrap.Response{
			Body:           nil,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}

// completePaymentPlanHandler completes a payment plan
// @Summary Completes a payment plan
// @Description completes a payment plan
// @Tags payment_plan
// @Produce json
// @Router /api/internal/pay_later/user/{user_uuid}/payment_plans/{uuid}/complete [post]
// @Param user_uuid path string true "User UUID"
// @Param uuid path string true "Payment Plan UUID"
// @Success 200
// @Failure 400 {object} handlerwrap.ErrorResponse "bad reqBody"
// @Failure 401 {object} handlerwrap.ErrorResponse "payment plan not belongs to user"
// @Failure 403 {object} handlerwrap.ErrorResponse "payment plan is not in pending"
// @Failure 404 {object} handlerwrap.ErrorResponse "payment plan not found"
// @Failure 500 {object} handlerwrap.ErrorResponse "internal error"
func completePaymentPlanHandler(paramsGetter handlerwrap.NamedURLParamsGetter) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		_, err := parsePaymentPlanParam(req.Context(), paramsGetter)
		if err != nil {
			return nil, err
		}

		return &handlerwrap.Response{
			Body:           nil,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}

type RefundRequest struct {
	RefundID   string      `json:"refund_id"`
	PaymentID  string      `json:"payment_id"`
	Amount     string      `json:"amount"`
	Currency   string      `json:"currency"`
	RefundData interface{} `json:"refund_data,omitempty"`
}

// refundHandler refunds a payment
// @Summary Refunds a payment
// @Description refunds a payment
// @Tags payment
// @Produce json
// @Router /api/internal/pay_later/refund [post]
// @Param refund_request body RefundRequest true "Refund reqBody data"
// @Success 200
// @Failure 403 {object} handlerwrap.ErrorResponse "payment plan is unconfirmed"
// @Failure 404 {object} handlerwrap.ErrorResponse "payment plan not found"
// @Failure 500 {object} handlerwrap.ErrorResponse "internal error"
func refundHandler() handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		var request RefundRequest

		body, readErr := io.ReadAll(req.Body)
		if readErr != nil {
			return nil, InvalidRequestBodyError{Err: readErr}.ToErrorResponse()
		}

		if err := json.Unmarshal(body, &request); err != nil {
			return nil, InvalidRequestBodyError{Err: err, Data: string(body)}.ToErrorResponse()
		}

		return &handlerwrap.Response{
			Body:           nil,
			HTTPStatusCode: http.StatusOK,
		}, nil
	}
}
