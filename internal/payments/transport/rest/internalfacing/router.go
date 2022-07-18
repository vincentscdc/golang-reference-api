package internalfacing

import (
	"golangreferenceapi/internal/payments/service"

	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
	"github.com/rs/zerolog"
)

func AddRoutes(
	router chi.Router, log *zerolog.Logger,
	paramsGetter handlerwrap.NamedURLParamsGetter,
	paymentService service.PaymentPlanService,
	version string,
) {
	router.Route("/internal/"+version, func(rtr chi.Router) {
		rtr.Post("/payment-plans",
			handlerwrap.Wrapper(log, createPendingPaymentPlanHandler(paymentService)))
		rtr.Post("/payment-plans/{payment_uuid}/complete",
			handlerwrap.Wrapper(log, completePaymentPlanHandler(paramsGetter, paymentService)))
	})
}
