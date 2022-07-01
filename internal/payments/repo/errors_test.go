package repo

import "testing"

func TestErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{
			name: "query run error",
			err:  pgxDBQueryRunError{},
			msg:  "failed to run query",
		},
		{
			name: "unsupported db entity",
			err:  unsupportedDBEntityError{},
			msg:  "DB entity interface does not match any supported struct",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err.Error() != tt.msg {
				t.Errorf("unexpected error, expected: %v, actual: %v", tt.msg, tt.err.Error())
			}
		})
	}
}
