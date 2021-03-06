package api

import (
	"fmt"
	"net/http"
	"os"

	"golangreferenceapi/internal/payments/docs"
	"golangreferenceapi/internal/payments/repo"
	"golangreferenceapi/internal/payments/service"
	"golangreferenceapi/internal/payments/transport/grpc/bnplapi/creditline/v1"
	grpcuserfacing "golangreferenceapi/internal/payments/transport/grpc/userfacing"
	"golangreferenceapi/internal/payments/transport/rest"
	"golangreferenceapi/internal/payments/transport/rest/internalfacing"
	"golangreferenceapi/internal/payments/transport/rest/userfacing"

	"github.com/monacohq/golang-common/transport/http/middleware/requestlogger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
)

func (s *API) setupLog() {
	if s.cfg.Application.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func (s *API) setupHTTPServer(repository repo.Repository) {
	// main router
	httpRouter := chi.NewRouter()
	httpRouter.Use(requestlogger.RequestLogger(&log.Logger))

	httpRouter.Mount("/debug", middleware.Profiler())

	httpRouter.Get("/sys/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	httpRouter.Get("/"+s.cfg.Application.Version+"/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("%s/%s/%s/swagger/doc.json",
				s.cfg.Application.URL.Schemes[0], s.cfg.Application.URL.Host, s.cfg.Application.Version,
			),
		),
	))

	paymentService := service.NewPaymentPlanService()
	paymentService.UseRepo(repository)

	httpRouter.Route("/", func(r chi.Router) {
		userfacing.AddRoutes(r, &log.Logger, paymentService, s.cfg.Application.Version)
		internalfacing.AddRoutes(r, &log.Logger, rest.ChiNamedURLParamsGetter, paymentService, s.cfg.Application.Version)
	})

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.Application.Port),
		ReadTimeout:       s.cfg.Application.Timeouts.ReadTimeout,
		ReadHeaderTimeout: s.cfg.Application.Timeouts.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.Application.Timeouts.WriteTimeout,
		IdleTimeout:       s.cfg.Application.Timeouts.IdleTimeout,
		Handler:           otelhttp.NewHandler(httpRouter, "server"),
	}
}

func (s *API) setupGRPCServer() {
	// grpc
	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	payLaterServer := grpcuserfacing.NewPayLaterServer()
	creditline.RegisterPayLaterServiceServer(s.grpcServer, payLaterServer)
}

func (s *API) setupSwagger() {
	// swagger
	version := "v1"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/" + version
	docs.SwaggerInfo.Host = s.cfg.Application.URL.Host
	docs.SwaggerInfo.Schemes = s.cfg.Application.URL.Schemes
}
