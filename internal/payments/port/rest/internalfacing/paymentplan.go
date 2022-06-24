package internalfacing

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golangreferenceapi/internal/payments/service"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

type CreatePendingPaymentPlanRequest struct {
	PendingPayment CreatePendingPayment `json:"payment"`
}

type CreatePendingPayment struct {
	ID           string              `json:"id"`
	Currency     string              `json:"currency"`
	TotalAmount  string              `json:"total_amount"`
	Installments []CreateInstallment `json:"installments"`
}

type CreateInstallment struct {
	ID       string `json:"id"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	DueAt    string `json:"due_at"`
}

type CreatePendingPaymentPlanResponse struct {
	PendingPayment service.PaymentPlans `json:"payment"`
}

type ListPaymentPlanResponse struct {
	Payment service.PaymentPlans `json:"payment"`
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
func createPendingPaymentPlanHandler(
	paramsGetter handlerwrap.NamedURLParamsGetter,
	paymentService service.PaymentPlanService,
) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		var (
			userUUID *uuid.UUID
			request  CreatePendingPaymentPlanRequest
			respErr  *handlerwrap.ErrorResponse
		)

		userUUID, respErr = parseUUIDFormatParam(req.Context(), paramsGetter, urlParamUserUUID)
		if respErr != nil {
			return nil, respErr
		}

		body, readErr := io.ReadAll(req.Body)
		if readErr != nil {
			return nil, InvalidRequestBodyError{Err: readErr}.ToErrorResponse()
		}

		if err := json.Unmarshal(body, &request); err != nil {
			return nil, InvalidRequestBodyError{Err: err, Data: string(body)}.ToErrorResponse()
		}

		createPaymentParams, err := NewCreatePaymentPlanParams(&request)
		if err != nil {
			return nil, InvalidRequestBodyError{Err: err, Data: string(body)}.ToErrorResponse()
		}

		paymentPlan, serviceErr := paymentService.CreatePendingPaymentPlan(req.Context(), *userUUID, createPaymentParams)
		if serviceErr != nil {
			return nil, InternalError{Err: serviceErr, Data: string(body)}.ToErrorResponse()
		}

		resp := CreatePendingPaymentPlanResponse{
			PendingPayment: *paymentPlan,
		}

		return &handlerwrap.Response{
			Body:           resp,
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
func completePaymentPlanHandler(
	paramsGetter handlerwrap.NamedURLParamsGetter,
	paymentService service.PaymentPlanService,
) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		var (
			paymentUUID   *uuid.UUID
			installmentID *uuid.UUID
			userUUID      *uuid.UUID
			respErr       *handlerwrap.ErrorResponse
		)

		userUUID, respErr = parseUUIDFormatParam(req.Context(), paramsGetter, urlParamUserUUID)
		if respErr != nil {
			return nil, respErr
		}

		paymentUUID, respErr = parseUUIDFormatParam(req.Context(), paramsGetter, urlParamPaymentUUID)
		if respErr != nil {
			return nil, respErr
		}

		installmentID, respErr = parseUUIDFormatParam(req.Context(), paramsGetter, urlParamInstallmentID)
		if respErr != nil {
			return nil, respErr
		}

		err := paymentService.CompletePaymentPlanCreation(context.Background(), *userUUID, *paymentUUID)
		if err != nil {
			data := fmt.Sprintf("err:invalid params,paymentUUID:%s,installmenetsID:%s", paymentUUID, installmentID)

			return nil, InternalError{Err: err, Data: data}.ToErrorResponse()
		}

		// optimize me :-(
		paymentPlan, serviceErr := paymentService.GetPaymentPlanByUserID(req.Context(), *userUUID)
		if serviceErr != nil {
			return nil, InternalError{Err: serviceErr, Data: "query paymentPlan error"}.ToErrorResponse()
		}

		resp := make([]ListPaymentPlanResponse, 0, len(paymentPlan))

		for _, plan := range paymentPlan {
			resp = append(resp, ListPaymentPlanResponse{Payment: plan})
		}

		return &handlerwrap.Response{
			Body:           resp,
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
