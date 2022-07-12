package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golangreferenceapi/internal/api/configuration"
	"golangreferenceapi/internal/payments/mock/repomock"
)

func Test_Start(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg := configuration.Config{}
	cfg.Application.Version = "v1"
	cfg.Application.Port = 8000
	cfg.Application.PrettyLog = true
	cfg.Application.URL.Host = "locahost:8000"
	cfg.Application.URL.Schemes = []string{"https"}
	cfg.Application.Timeouts.ReadTimeout = 2 * time.Second
	cfg.Application.Timeouts.ReadHeaderTimeout = 1 * time.Second
	cfg.Application.Timeouts.WriteTimeout = 2 * time.Second
	cfg.Application.Timeouts.IdleTimeout = 1 * time.Minute
	cfg.Grpc.Port = 9000
	cfg.Observability.Collector.Host = "opentelemetry-collector.otel-collector"
	cfg.Observability.Collector.Port = 4317

	apiSrv := NewAPI(&cfg, &repomock.MockRepository{})

	// append mock err to test handling of shutdownFuncs which return err
	apiSrv.shutdownFuncs = append(apiSrv.shutdownFuncs, &shutdownFunc{
		f: func() error {
			return fmt.Errorf("mock err")
		},
		msg: "mock err",
	})

	shutdown, err := apiSrv.Start(ctx)
	if err != nil {
		t.Errorf("api failed to start: %v", err)
	}

	shutdown()
}
