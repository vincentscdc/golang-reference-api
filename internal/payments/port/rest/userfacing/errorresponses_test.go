package userfacing

import (
	"errors"
	"testing"
)

func TestInvalidRequestBodyError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Data string
		Err  error
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				Data: `{"x}`,
				Err:  errors.New("invalid format"),
			},
			want: "can not unmarshal request body `{\"x}`: invalid format",
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

			if got := err.Error(); got != tt.want {
				t.Errorf("returned unexpected error message got %v want %v", got, tt.want)
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
		name   string
		fields fields
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				Data: `{}`,
				Err:  errors.New("some error occur"),
			},
			want: "internal error `{}`: some error occur",
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
		})
	}
}
