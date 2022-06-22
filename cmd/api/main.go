package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golangreferenceapi/internal/configuration"
	"golangreferenceapi/internal/docs"
	"golangreferenceapi/internal/port/grpc/bnplapi/creditline/v1"
	grpcuserfacing "golangreferenceapi/internal/port/grpc/userfacing"
	"golangreferenceapi/internal/port/rest"
	"golangreferenceapi/internal/port/rest/internalfacing"
	"golangreferenceapi/internal/port/rest/userfacing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/monacohq/golang-common/monitoring/otelinit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
)

// nolint: gochecknoglobals // only allowed global vars - filled at build time - do not change
var (
	CommitTime = "dev"
	CommitHash = "dev"
)

// @title bnpl API
// @description This is the service used for buy now pay later.
// @version

// @contact.name Vincent
// @contact.email vincent.serpoul@crypto.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() { // nolint: cyclop // temporary, will be moved to multiple funcs
	// context
	ctx := context.Background()

	// configuration
	currEnv := "local"
	if e := os.Getenv("APP_ENVIRONMENT"); e != "" {
		currEnv = e
	}

	cfg, err := configuration.GetConfig(currEnv)
	if err != nil {
		if errors.Is(err, configuration.MissingBaseConfigError{}) {
			log.Printf("getConfig: %v", err)

			return
		}

		log.Printf("getConfig: %v", err)
	}

	// logging
	if cfg.Application.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// otel
	shutdown, err := otelinit.InitProvider(
		ctx,
		"bnpl",
		otelinit.WithGRPCTraceExporter(
			ctx,
			fmt.Sprintf("%s:%d", cfg.Observability.Collector.Host, cfg.Observability.Collector.Port),
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize opentelemetry")

		return
	}

	defer func() {
		if errShutdown := shutdown(); errShutdown != nil {
			log.Warn().Err(errShutdown).Msg("shutdown")
		}
	}()

	// main router
	httpRouter := chi.NewRouter()

	httpRouter.Mount("/debug", middleware.Profiler())

	httpRouter.Get("/sys/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	version := "v1"

	// swagger
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/" + version
	docs.SwaggerInfo.Host = cfg.Application.URL.Host
	docs.SwaggerInfo.Schemes = cfg.Application.URL.Schemes

	httpRouter.Get("/"+version+"/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(
			fmt.Sprintf("%s/%s/%s/swagger/doc.json",
				cfg.Application.URL.Schemes[0], cfg.Application.URL.Host, version,
			),
		),
	))

	httpRouter.Route("/"+version, func(r chi.Router) {
		userfacing.AddRoutes(r, &log.Logger)
		internalfacing.AddRoutes(r, &log.Logger, rest.ChiNamedURLParamsGetter)
	})

	// serve router
	log.Info().
		Int("port", cfg.Application.Port).
		Str("commit time", CommitTime).
		Str("commit hash", CommitHash).
		Msg("listening")

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Application.Port),
		ReadTimeout:       cfg.Application.Timeouts.ReadTimeout,
		ReadHeaderTimeout: cfg.Application.Timeouts.ReadHeaderTimeout,
		WriteTimeout:      cfg.Application.Timeouts.WriteTimeout,
		IdleTimeout:       cfg.Application.Timeouts.IdleTimeout,
		Handler:           otelhttp.NewHandler(httpRouter, "server"),
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start http server
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("")
		}

		// Wait for server context to be stopped
		<-serverCtx.Done()
	}()

	// start grpc server
	addr := fmt.Sprintf(":%d", cfg.Grpc.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msgf("failed to listen on address %s", addr)

		return
	}

	grpcServer := grpc.NewServer()
	server := grpcuserfacing.NewPayLaterServer()

	creditline.RegisterPayLaterServiceServer(grpcServer, server)

	go func() {
		log.Info().Msgf("start to listen grpc server on %s", addr)

		err = grpcServer.Serve(lis)
		if err != nil {
			log.Error().Err(err).Msg("failed to serve grpc server")

			return
		}
	}()

	<-sig

	const shutdownGracePeriod = 30 * time.Second

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, cancel := context.WithTimeout(serverCtx, shutdownGracePeriod)
	defer cancel()

	go func() {
		<-shutdownCtx.Done()

		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			log.Fatal().Err(err).Msg("graceful shutdown timed out.. forcing exit.")
		}
	}()

	// Trigger graceful shutdown
	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Err(err).Msg("error shutting down")
	}

	grpcServer.GracefulStop()

	serverStopCtx()
}
