package userfacing

import (
	"net/http"

	"golangreferenceapi/internal/payments/service"
	"golangreferenceapi/internal/payments/transport/rest"

	"github.com/monacohq/golang-common/transport/http/middleware/cryptouseruuid"

	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
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

// listPaymentPlansHandler renders the payment plans
// @Summary Renders a user's payment plans
// @Description returns pagination of one user's payment plans
// @Tags payment_plan
// @Produce json
// @Router /api/v1/payment-plans [get]
// @Param X-CRYPTO-USER-UUID header string true "User UUID"
// @Param offset query int64 false "Start index in the list" minimum(0)
// @Param limit query int64 false "Number of items displayed" minimum(0) maximum(10)
// @Param created_at_order query string false "Order by payment.created_at asc  OR desc" Enums(asc, desc) default(desc)
// @Success 200 {object} PaymentPlanResponse
func listPaymentPlansHandler(paymentService service.PaymentPlanService) handlerwrap.TypedHandler {
	return func(req *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		// user uuid
		uid, err := cryptouseruuid.GetUserUUID(req.Context())
		if err != nil {
			return nil, handlerwrap.NewErrorResponseFromCryptoUserUUIDError(err)
		}
		// pagination params
		_, paginateErr := rest.ParsePaginationURLQuery(req.URL, paymentPlansDefaultLimit, rest.PaymentPlansCreatedAtOrderDESC)
		if paginateErr != nil {
			return nil, paginateErr
		}

		paymentPlans, serviceErr := paymentService.GetPaymentPlanByUserID(req.Context(), *uid)
		if serviceErr != nil {
			return nil, rest.ServiceErrorToErrorResp(serviceErr)
		}

		resp := &handlerwrap.Response{Body: PaymentPlanResponse{Payments: paymentPlans}, StatusCode: http.StatusOK}

		return resp, nil
	}
}
