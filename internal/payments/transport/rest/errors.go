package rest

import (
	"errors"
	"net/http"

	"golangreferenceapi/internal/payments/service"

	"github.com/monacohq/golang-common/transport/http/handlerwrap/v3"
)

func ServiceErrorToErrorResp(err error) *handlerwrap.ErrorResponse {
	switch {
	case errors.Is(err, service.ErrRecordNotFound):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusNotFound,
			"record_not_found",
			"record not found",
		)
	case errors.Is(err, service.ErrGenerateUUID):
		return handlerwrap.NewErrorResponse(
			err,
			make(map[string]string),
			http.StatusInternalServerError,
			"uuid_generate_failed",
			"uuid generated failed",
		)
	default:
		return handlerwrap.InternalServerError{Err: err}.ToErrorResponse()
	}
}
