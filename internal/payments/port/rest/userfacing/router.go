package userfacing

import (
	"golangreferenceapi/internal/payments/service"

	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
	"github.com/monacohq/golang-common/transport/http/middleware/cryptouseruuid"
	"github.com/rs/zerolog"
)

const (
	responseKeyPaymentPlans = "payment_plans"
)

func AddRoutes(
	router chi.Router, log *zerolog.Logger, paymentService service.PaymentPlanService,
) {
	router.Route("/api/pay_later/", func(r chi.Router) {
		r.Use(cryptouseruuid.UserUUID(log))
		r.Get("/payment_plans", handlerwrap.Wrapper(log, listPaymentPlansHandlerOKStyle(log, paymentService)))
	})
}

// PaymentPlansResponseOKStyle represents the response struct
type PaymentPlansResponseOKStyle struct {
	OKStyleResponseBase
	PaymentPlans []PaymentPlanResponse `json:"payment_plans,omitempty"`
}

// listPaymentPlansHandlerOKStyle renders the payment plans
// @Summary Renders a user's payment plans
// @Description returns pagination of one user's payment plans
// @Tags payment_plan
// @Produce json
// @Router /api/pay_later/payment_plans [get]
// @Param X-CRYPTO-USER-UUID header string true "User UUID"
// @Param offset query int64 false "Start index in the list" minimum(0)
// @Param limit query int64 false "Number of items displayed" minimum(0) maximum(10)
// @Param created_at_order query string false "Order by payment.created_at asc  OR desc" Enums(asc, desc) default(desc)
// @Success 200 {object} PaymentPlansResponseOKStyle
func listPaymentPlansHandlerOKStyle(
	log *zerolog.Logger,
	paymentService service.PaymentPlanService,
) handlerwrap.TypedHandler {
	return OKStyleWrapper(log, responseKeyPaymentPlans, listPaymentPlansHandler(paymentService))
}
