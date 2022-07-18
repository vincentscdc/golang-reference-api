package internalfacing

import (
	"net/http"

	"golangreferenceapi/internal/payments/service"
	"golangreferenceapi/internal/payments/transport/rest"

	"github.com/gofrs/uuid"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
)

type CreatePendingPaymentPlanRequest struct {
	PendingPayment service.CreatePaymentPlanParams `json:"payment"`
}

type CreatePendingPaymentPlanResponse struct {
	PendingPayment service.PaymentPlans `json:"payment"`
}
type CompletePaymentPlanRequest struct {
	Payment service.CompletePaymentPlanParams `json:"payment"`
}

type CompletePaymentPlanResponse struct {
	Payment service.PaymentPlans `json:"payment"`
}

type ListPaymentPlanResponse struct {
	Payment service.PaymentPlans `json:"payment"`
}

// createPendingPaymentPlanHandler creates a pending payment plan
// @Summary Creates a pending a payment plan
// @Description pre creates a payment plan
// @Tags payment_plan
// @Produce json
// @Router /internal/v1/payment_plans [post]
// @Param pre_create_payment_plan_request body CreatePendingPaymentPlanRequest true "Pre create payment plan reqBody"
// @Success 200 {object} CreatePendingPaymentPlanResponse
// @Failure 400 {object} handlerwrap.ErrorResponse "bad reqBody"
// @Failure 500 {object} handlerwrap.ErrorResponse "internal error"
func createPendingPaymentPlanHandler(
	paymentService service.PaymentPlanService,
) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		var request CreatePendingPaymentPlanRequest

		if errResp := handlerwrap.BindBody(req, &request); errResp != nil {
			return nil, errResp
		}

		paymentPlan, err := paymentService.CreatePendingPaymentPlan(req.Context(), &request.PendingPayment)
		if err != nil {
			return nil, rest.ServiceErrorToErrorResp(err)
		}

		resp := CreatePendingPaymentPlanResponse{
			PendingPayment: *paymentPlan,
		}

		return &handlerwrap.Response{
			Body:       resp,
			StatusCode: http.StatusOK,
		}, nil
	}
}

// completePaymentPlanHandler completes a payment plan
// @Summary Completes a payment plan
// @Description completes a payment plan
// @Tags payment_plan
// @Produce json
// @Router /internal/v1/payment-plans/{uuid}/complete [post]
// @Param complete_payment_plan_request body CompletePaymentPlanRequest true "Complete payment plan reqBody"
// @Param uuid path string true "Payment Plan UUID"
// @Success 200 {object} CompletePaymentPlanResponse
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
			paymentUUID *uuid.UUID
			respErr     *handlerwrap.ErrorResponse
			request     CompletePaymentPlanRequest
		)

		if errResp := handlerwrap.BindBody(req, &request); errResp != nil {
			return nil, errResp
		}

		paymentUUID, respErr = parseUUIDFormatParam(req.Context(), paramsGetter, urlParamPaymentUUID)
		if respErr != nil {
			return nil, respErr
		}

		payment, err := paymentService.CompletePaymentPlanCreation(req.Context(), *paymentUUID, &request.Payment)
		if err != nil {
			return nil, rest.ServiceErrorToErrorResp(err)
		}

		return &handlerwrap.Response{
			Body:       CompletePaymentPlanResponse{Payment: *payment},
			StatusCode: http.StatusOK,
		}, nil
	}
}
