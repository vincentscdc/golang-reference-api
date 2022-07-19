package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golangreferenceapi/internal/api/configuration"
	"golangreferenceapi/internal/payments/repo"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

//nolint: gochecknoglobals // only allowed global vars - filled at build time - do not change
var (
	CommitTime = "dev"
	CommitHash = "dev"
)

// graceful shutdown functionality can be implemented as a separate common library in future
// https://github.com/monacohq/golang-reference-api/pull/26#issuecomment-1162779857
type shutdownFunc struct {
	f   func() error
	msg string
}

type API struct {
	httpServer    *http.Server
	grpcServer    *grpc.Server
	cfg           configuration.Config
	shutdownFuncs []*shutdownFunc
}

func NewAPI(cfg *configuration.Config, repository repo.Repository) *API {
	srv := &API{cfg: *cfg}
	srv.setupLog()
	srv.setupHTTPServer(repository)
	srv.setupGRPCServer()
	srv.setupSwagger()

	return srv
}

func (s *API) Start(ctx context.Context) (func(), error) {
	log.Info().
		Str("commit time", CommitTime).
		Str("commit hash", CommitHash).
		Msg("listening")

	if err := s.startOtel(ctx); err != nil {
		return nil, fmt.Errorf("failed to start otel: %w", err)
	}

	serverCtx, serverStopCtx := s.startHTTPServer()

	if err := s.startGRPCServer(); err != nil {
		return nil, fmt.Errorf("failed to start grpc server: %w", err)
	}

	s.shutdownFuncs = append(s.shutdownFuncs, &shutdownFunc{
		f: func() error {
			serverStopCtx()

			return nil
		},
		msg: "serverStopCtx",
	})

	return func() {
		s.shutdown(serverCtx)
	}, nil
}

func (s *API) shutdown(serverCtx context.Context) {
	for _, sdf := range s.shutdownFuncs {
		if errShutdown := sdf.f(); errShutdown != nil {
			log.Warn().Err(errShutdown).Msg(sdf.msg)
		}
	}

	const shutdownGracePeriod = 30 * time.Second

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, cancel := context.WithTimeout(serverCtx, shutdownGracePeriod)
	defer cancel()

	go func() {
		<-shutdownCtx.Done()

		if err := shutdownCtx.Err(); errors.Is(err, context.DeadlineExceeded) {
			log.Fatal().Err(err).Msg("graceful shutdown timed out.. forcing exit.")
		}
	}()

	// Trigger graceful shutdown
	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Error().Err(err).Msg("error shutting down")
	}

	s.grpcServer.GracefulStop()
}
