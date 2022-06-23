package internalfacing

import (
	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
	"github.com/rs/zerolog"

	"golangreferenceapi/internal/payments/service"
)

func AddRoutes(
	router chi.Router, log *zerolog.Logger,
	paramsGetter handlerwrap.NamedURLParamsGetter,
	paymentService service.PaymentPlanService,
) {
	router.Route("/api/internal/pay_later/", func(rtr chi.Router) {
		rtr.Get("/user/{user_uuid}/credit_line",
			handlerwrap.Wrapper(log, getCreditLineHandler(paramsGetter)))
		rtr.Post("/user/{user_uuid}/payment_plans",
			handlerwrap.Wrapper(log, createPendingPaymentPlanHandler(paramsGetter, paymentService)))
		rtr.Post("/user/{user_uuid}/payment_plans/{uuid}/cancel",
			handlerwrap.Wrapper(log, cancelPaymentPlanHandler(paramsGetter)))
		rtr.Post("/user/{user_uuid}/payment_plans/{payment_uuid}/installments/{installments_id}/payment",
			handlerwrap.Wrapper(log, completePaymentPlanHandler(paramsGetter, paymentService)))
		rtr.Post("/refund", handlerwrap.Wrapper(log, refundHandler()))
	})
}
