package userfacing

import (
	"fmt"
	"net/http"

	"github.com/monacohq/golang-common/transport/http/handlerwrap"
)

type InvalidRequestBodyError struct {
	Data string
	Err  error
}

func (err InvalidRequestBodyError) Error() string {
	return fmt.Sprintf("can not unmarshal request body `%v`: %v", err.Data, err.Err)
}

func (err InvalidRequestBodyError) ToErrorResponse() *handlerwrap.ErrorResponse {
	return &handlerwrap.ErrorResponse{
		Error:          err.Err,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorCode:      "invalid_request_body",
		ErrorMsg:       err.Error(),
	}
}

type InternalError struct {
	Data string
	Err  error
}

func (err InternalError) Error() string {
	return fmt.Sprintf("internal error `%v`: %v", err.Data, err.Err)
}

func (err InternalError) ToErrorResponse() *handlerwrap.ErrorResponse {
	return &handlerwrap.ErrorResponse{
		Error:          err.Err,
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorCode:      "user_error",
		ErrorMsg:       err.Error(),
	}
}
