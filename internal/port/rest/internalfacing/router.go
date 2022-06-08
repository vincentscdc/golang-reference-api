package internalfacing

import (
	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap"
	"github.com/rs/zerolog"
)

func AddRoutes(
	router chi.Router, log *zerolog.Logger,
	paramsGetter handlerwrap.NamedURLParamsGetter,
) {
	router.Route("/api/internal/pay_later/", func(rtr chi.Router) {
		rtr.Get("/user/{user_uuid}/credit_line",
			handlerwrap.Wrapper(log, getCreditLineHandler(paramsGetter)))
		rtr.Post("/user/{user_uuid}/payment_plans",
			handlerwrap.Wrapper(log, createPendingPaymentPlanHandler(paramsGetter)))
		rtr.Post("/user/{user_uuid}/payment_plans/{uuid}/cancel",
			handlerwrap.Wrapper(log, cancelPaymentPlanHandler(paramsGetter)))
		rtr.Post("/user/{user_uuid}/payment_plans/{uuid}/complete",
			handlerwrap.Wrapper(log, completePaymentPlanHandler(paramsGetter)))
		rtr.Post("/refund", handlerwrap.Wrapper(log, refundHandler()))
	})
}
