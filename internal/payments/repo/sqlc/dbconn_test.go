package sqlc

import (
	"context"
	"strings"
	"testing"
	"time"

	"golangreferenceapi/internal/api/configuration"
)

func Test_NewRepo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name    string
		cfg     *configuration.Database
		wantErr bool
	}{
		{
			name: "happy path",
			cfg: &configuration.Database{
				Host:         "localhost",
				Port:         strings.Split(getHostPort(testRefDockertestResource, "5432/tcp"), ":")[1],
				User:         "postgres",
				Password:     "postgres",
				Database:     "datawarehouse",
				MaxConns:     10,
				MaxIdleConns: 10,
				MaxLifeTime:  1 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "unhappy path - fail to init PGInit",
			cfg: &configuration.Database{
				Host:         "wrong:localhost",
				Port:         "wrong/port",
				User:         "postgres",
				Password:     "postgres",
				Database:     "datawarehouse",
				MaxConns:     10,
				MaxIdleConns: 10,
				MaxLifeTime:  1 * time.Minute,
			},
			wantErr: true,
		},
		{
			name:    "unhappy path - fail to init conn pool",
			cfg:     &configuration.Database{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewRepo(ctx, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected err result: %v", err)
			}
		})
	}
}
