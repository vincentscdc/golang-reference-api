package rest

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestChiNamedURLParamsGetter(t *testing.T) {
	t.Parallel()

	var (
		key   = "key"
		value = "value"
	)

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add(key, value)

	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routerCtx)

	v, err := ChiNamedURLParamsGetter(ctx, key)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if v != value {
		t.Errorf("expected: %v, actual: %v", value, v)
	}

	_, err = ChiNamedURLParamsGetter(ctx, "unknown_key")
	if err == nil {
		t.Errorf("unexpected nil error")
	} else if err.StatusCode != http.StatusBadRequest {
		t.Errorf(
			"expected: %v, actual: %v", err.StatusCode, http.StatusBadRequest,
		)
	}
}
