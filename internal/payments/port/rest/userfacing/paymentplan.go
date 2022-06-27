package userfacing

import (
	"context"
	"net/http"

	"golangreferenceapi/internal/payments/port/rest"

	"golangreferenceapi/internal/payments/service"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

const (
	paymentPlansDefaultLimit = 10
)

// PaymentPlanResponse represents a specific payment plan
type PaymentPlanResponse struct {
	Payments []service.PaymentPlans `json:"payments"`
}

type Installments struct {
	ID       string `json:"id"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	DueAt    string `json:"due_at"`
	Status   string `json:"status"`
}

func listPaymentPlansHandler(paymentService service.PaymentPlanService) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		// user uuid
		uid, err := getUserUUID(req.Context())
		if err != nil {
			return nil, err
		}
		// pagination params
		_, err = parsePaginationURLQuery(req.URL, paymentPlansDefaultLimit, paymentPlansCreatedAtOrderDESC)
		if err != nil {
			return nil, err
		}

		paymentPlans, serviceErr := paymentService.GetPaymentPlanByUserID(context.Background(), *uid)
		if serviceErr != nil {
			return nil, rest.ServiceErrorToErrorResp(serviceErr)
		}

		resp := &handlerwrap.Response{Body: PaymentPlanResponse{Payments: paymentPlans}, HTTPStatusCode: http.StatusOK}

		return resp, nil
	}
}
