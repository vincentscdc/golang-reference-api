package userfacing

import (
	"golangreferenceapi/internal/payments/service"

	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v3"
	"github.com/monacohq/golang-common/transport/http/middleware/cryptouseruuid"
	"github.com/rs/zerolog"
)

func AddRoutes(
	router chi.Router,
	log *zerolog.Logger,
	paymentService service.PaymentPlanService,
	version string,
) {
	router.Route("/api/"+version, func(r chi.Router) {
		r.Use(cryptouseruuid.UserUUID(log))
		r.Get("/payment-plans", handlerwrap.Wrapper(log, listPaymentPlansHandler(paymentService)))
	})
}
