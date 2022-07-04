package api

import (
	"context"
	"errors"
	"testing"

	"golangreferenceapi/internal/api/configuration"
)

func Test_Start(t *testing.T) {
	t.Parallel()

	cfg, err := configuration.GetConfig("../../config/api", "localdev")
	if err != nil && errors.Is(err, configuration.MissingBaseConfigError{}) {
		t.Fatal(err)
	}

	ctx := context.Background()

	apiSrv, err := NewAPI(ctx, &cfg)
	if err != nil {
		t.Errorf("api failed to initialize: %v", err)
	}

	shutdown, err := apiSrv.Start(ctx)
	if err != nil {
		t.Errorf("api failed to start: %v", err)
	}

	shutdown()
}
