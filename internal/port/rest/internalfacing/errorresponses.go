package internalfacing

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
