package configuration

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

func Test_GetConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		env            string
		expectedErr    interface{}
		expectedVerRgx string
	}{
		{
			name:           "happy path",
			path:           "../../../config/api",
			env:            "base",
			expectedErr:    nil,
			expectedVerRgx: "v[0-9]*\\.*[0-9]*\\.*[0-9]*",
		},
		{
			name:           "unhappy path - missing env config",
			path:           "../../../config/api",
			env:            "",
			expectedErr:    &MissingEnvConfigError{},
			expectedVerRgx: "v[0-9]*\\.*[0-9]*\\.*[0-9]*",
		},
		{
			name:           "unhappy path - missing base config",
			path:           "./wrongpath",
			env:            "",
			expectedErr:    &MissingBaseConfigError{},
			expectedVerRgx: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg, err := GetConfig(tt.path, tt.env)

			if tt.expectedErr != nil && !errors.As(err, tt.expectedErr) {
				t.Errorf("unexpected err received: %v", err)
			}

			match, _ := regexp.MatchString(tt.expectedVerRgx, cfg.Application.Version)
			if !match {
				t.Errorf("unexpected version value: %s", cfg.Application.Version)
			}
		})
	}
}

func Test_MissingEnvConfigError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            MissingEnvConfigError{env: "local", err: fmt.Errorf("some error")},
			expectedString: "missing config local: some error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err.Error() != tt.expectedString {
				t.Errorf("unexpected Error string")
			}
		})
	}
}

func Test_AsMissingEnvConfigError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		err           error
		comparableErr interface{}
		expectedBool  bool
	}{
		{
			name:          "happy path",
			err:           MissingEnvConfigError{env: "local", err: fmt.Errorf("some error")},
			comparableErr: &MissingEnvConfigError{},
			expectedBool:  true,
		},
		{
			name:          "unhappy path",
			err:           MissingEnvConfigError{env: "local", err: fmt.Errorf("some error")},
			comparableErr: &MissingBaseConfigError{},
			expectedBool:  false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if errors.As(tt.err, tt.comparableErr) != tt.expectedBool {
				t.Errorf("unexpected bool")
			}
		})
	}
}

func Test_MissingBaseConfigError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedString string
	}{
		{
			name:           "happy path",
			err:            MissingBaseConfigError{err: fmt.Errorf("some error")},
			expectedString: "missing base config: some error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err.Error() != tt.expectedString {
				t.Errorf("unexpected Error string")
			}
		})
	}
}

func Test_AsMissingBaseConfigError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		err           error
		comparableErr interface{}
		expectedBool  bool
	}{
		{
			name:          "happy path",
			err:           MissingBaseConfigError{err: fmt.Errorf("some error")},
			comparableErr: &MissingBaseConfigError{},
			expectedBool:  true,
		},
		{
			name:          "unhappy path",
			err:           MissingBaseConfigError{err: fmt.Errorf("some error")},
			comparableErr: &MissingEnvConfigError{},
			expectedBool:  false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if errors.As(tt.err, tt.comparableErr) != tt.expectedBool {
				t.Errorf("unexpected bool")
			}
		})
	}
}
