package rest

import (
	"errors"
	"net/http"

	"golangreferenceapi/internal/payments/service"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

func ServiceErrorToErrorResp(err error) *handlerwrap.ErrorResponse {
	switch {
	case errors.Is(err, service.ErrRecordNotFound):
		return handlerwrap.NewErrorResponse(
			err,
			http.StatusNotFound,
			"record_not_found",
			"record not found",
		)
	default:
		return handlerwrap.InternalServerError{Err: err}.ToErrorResponse()
	}
}
