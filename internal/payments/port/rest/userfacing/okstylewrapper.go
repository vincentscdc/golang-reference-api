package userfacing

import (
	"net/http"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
	"github.com/rs/zerolog"
)

type OKStyleResponseBase struct {
	OK           bool   `json:"ok"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// OKStyleWrapper wraps a handler, returns a handler which returns a 200 *Response in ok style and nil *ErrorResponse.
// Log the error when handler returns a non-nil *ErrorResponse
func OKStyleWrapper(log *zerolog.Logger, bodyName string, handler handlerwrap.TypedHandler) handlerwrap.TypedHandler {
	return func(r *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
		resp, err := handler(r)
		if err == nil {
			return &handlerwrap.Response{
				Body: map[string]any{
					"ok":     true,
					bodyName: resp.Body,
				},
				HTTPStatusCode: http.StatusOK,
			}, nil
		}

		// ignore non-nil ErrorResponse and log error
		log.Error().
			Err(err.Error).
			Str("ErrorCode", err.ErrorCode).
			Int("HTTPStatusCode", err.HTTPStatusCode).
			Msg(err.ErrorMsg)

		response := &handlerwrap.Response{
			Body: map[string]any{
				"ok":            false,
				"error":         err.ErrorCode,
				"error_message": err.ErrorMsg,
			},
			HTTPStatusCode: http.StatusOK,
		}

		return response, nil
	}
}
