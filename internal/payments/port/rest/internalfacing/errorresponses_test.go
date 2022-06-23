package internalfacing

import (
	"errors"
	"net/http"
	"testing"
)

func TestInvalidRequestBodyError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Data string
		Err  error
	}

	tests := []struct {
		name       string
		fields     fields
		errStr     string
		statusCode int
	}{
		{
			name: "happy path",
			fields: fields{
				Data: `{"x}`,
				Err:  errors.New("invalid format"),
			},
			errStr:     "can not unmarshal request body `{\"x}`: invalid format",
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := InvalidRequestBodyError{
				Data: tt.fields.Data,
				Err:  tt.fields.Err,
			}

			if got := err.Error(); got != tt.errStr {
				t.Errorf("returned unexpected error message got %v want %v", got, tt.errStr)
			}

			errResponse := err.ToErrorResponse()
			if errResponse.HTTPStatusCode != tt.statusCode {
				t.Errorf("return unexpected status code got %v want %v", errResponse.HTTPStatusCode, tt.statusCode)
			}
		})
	}
}

func TestInternalError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Data string
		Err  error
	}

	tests := []struct {
		name       string
		fields     fields
		want       string
		statusCode int
	}{
		{
			name: "happy path",
			fields: fields{
				Data: `{}`,
				Err:  errors.New("some error occur"),
			},
			want:       "internal error `{}`: some error occur",
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := InternalError{
				Data: tt.fields.Data,
				Err:  tt.fields.Err,
			}

			if got := err.Error(); got != tt.want {
				t.Errorf("returned unexpected error message got %v want %v", got, tt.want)
			}

			errResponse := err.ToErrorResponse()
			if errResponse.HTTPStatusCode != tt.statusCode {
				t.Errorf("return unexpected status code got %v want %v", errResponse.HTTPStatusCode, tt.statusCode)
			}
		})
	}
}
