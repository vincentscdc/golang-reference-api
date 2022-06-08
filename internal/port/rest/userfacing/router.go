package userfacing

import (
	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
	"github.com/rs/zerolog"
)

const (
	responseKeyCreditInfo   = "credit_info"
	responseKeyPaymentPlans = "payment_plans"
)

func AddRoutes(
	router chi.Router, log *zerolog.Logger,
) {
	router.Route("/api/pay_later/", func(r chi.Router) {
		r.Use(UserUUID(log))
		r.Get("/credit_line", handlerwrap.Wrapper(log, getCreditLineHandlerOKStyle(log)))
		r.Get("/payment_plans", handlerwrap.Wrapper(log, listPaymentPlansHandlerOKStyle(log)))
	})
}

// CreditLineResponseOKStyle represents the response struct
type CreditLineResponseOKStyle struct {
	OKStyleResponseBase
	CreditLine *CreditLineResponse `json:"credit_info,omitempty"`
}

// getCreditLineHandler renders the credit line
// @Summary Renders a user's credit line info
// @Description returns the amount and status
// @Tags credit_line
// @Produce json
// @Router /api/pay_later/credit_line [get]
// @Param X-CRYPTO-USER-UUID header string true "User UUID"
// @Success 200 {object} CreditLineResponseOKStyle
func getCreditLineHandlerOKStyle(log *zerolog.Logger) handlerwrap.TypedHandler {
	return OKStyleWrapper(log, responseKeyCreditInfo, getCreditLineHandler())
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
func listPaymentPlansHandlerOKStyle(log *zerolog.Logger) handlerwrap.TypedHandler {
	return OKStyleWrapper(log, responseKeyPaymentPlans, listPaymentPlansHandler())
}
