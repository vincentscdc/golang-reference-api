package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/monacohq/golang-common/monitoring/otelinit"
	"github.com/rs/zerolog/log"
)

func (s *API) startOtel(ctx context.Context) error {
	var err error

	shutdownOtel, err := otelinit.InitProvider(
		ctx,
		"bnpl",
		otelinit.WithGRPCTraceExporter(
			ctx,
			fmt.Sprintf("%s:%d", s.cfg.Observability.Collector.Host, s.cfg.Observability.Collector.Port),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize opentelemetry: %w", err)
	}

	s.shutdownFuncs = append(s.shutdownFuncs, &shutdownFunc{f: shutdownOtel, msg: "shutdown otel"})

	return nil
}

func (s *API) startHTTPServer() (context.Context, context.CancelFunc) {
	// serve router
	log.Info().
		Int("http port", s.cfg.Application.Port).
		Str("host", s.cfg.Application.URL.Host).
		Str("readTimeout", s.cfg.Application.Timeouts.ReadTimeout.String()).
		Str("readHeaderTimeout", s.cfg.Application.Timeouts.ReadHeaderTimeout.String()).
		Str("writeTimeout", s.cfg.Application.Timeouts.WriteTimeout.String()).
		Str("idleTimeout", s.cfg.Application.Timeouts.IdleTimeout.String()).
		Msg("start HTTP server")

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Start http server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("")
		}

		// Wait for server context to be stopped
		<-serverCtx.Done()
	}()

	return serverCtx, serverStopCtx
}

func (s *API) startGRPCServer() error {
	// start grpc server
	addr := fmt.Sprintf(":%d", s.cfg.Grpc.Port)

	log.Info().
		Int("grpc port", s.cfg.Grpc.Port).
		Str("tcp", addr).
		Msg("start gRPC server")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %w", addr, err)
	}

	go func() {
		log.Info().Msgf("grpc server started listening on %s", addr)

		err = s.grpcServer.Serve(lis)
		if err != nil {
			log.Error().Err(err).Msg("failed to serve grpc server")

			return
		}
	}()

	return nil
}
