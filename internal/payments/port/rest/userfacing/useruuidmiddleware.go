package userfacing

import (
	"context"
	"errors"
	"net/http"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	HTTPHeaderKeyUserUUID = "X-CRYPTO-USER-UUID"
)

var ErrUserIDNotFound = errors.New("user id not found")

// UserUUID a middleware to get the user uuid from HTTP header.
// set it into ctx otherwise abort with a 401 HTTP status code
func UserUUID(log *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
			uuidVal := req.Header.Get(HTTPHeaderKeyUserUUID)
			if uuidVal == "" {
				respWriter.WriteHeader(http.StatusUnauthorized)

				log.Error().
					Str("path", req.URL.Path).
					Str("method", req.Method).
					Msg("uuid is empty")

				return
			}

			userUUID, err := uuid.Parse(uuidVal)
			if err != nil {
				log.Error().
					Err(err).
					Str("UserID", uuidVal).
					Msg("invalid user id")
				respWriter.WriteHeader(http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(respWriter, req.WithContext(setUserUUID(req.Context(), &userUUID)))
		})
	}
}

func setUserUUID(ctx context.Context, userUUID *uuid.UUID) context.Context {
	return context.WithValue(ctx, contextValKeyUserUUID, userUUID)
}

func getUserUUID(ctx context.Context) (*uuid.UUID, *handlerwrap.ErrorResponse) {
	userUUID, ok := ctx.Value(contextValKeyUserUUID).(*uuid.UUID)
	if ok {
		return userUUID, nil
	}

	return nil,
		handlerwrap.NewErrorResponse(ErrUserIDNotFound, http.StatusUnauthorized, "", "user id not found")
}
