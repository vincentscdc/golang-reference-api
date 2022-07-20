package rest

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/monacohq/golang-common/transport/http/handlerwrap/v3"
)

func ChiNamedURLParamsGetter(ctx context.Context, key string) (string, *handlerwrap.ErrorResponse) {
	v := chi.URLParamFromCtx(ctx, key)
	if v == "" {
		return "", handlerwrap.MissingParamError{Name: key}.ToErrorResponse()
	}

	return v, nil
}
