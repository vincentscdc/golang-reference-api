package userfacing

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_getPaginationURLQuery(t *testing.T) {
	t.Parallel()

	type args struct {
		rawURL                string
		defaultLimit          int64
		defaultCreatedAtOrder string
	}

	tests := []struct {
		name    string
		args    args
		want    *PaginationURLQuery
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				rawURL:                "https://example.com?offset=0&limit=10&created_at_order=asc",
				defaultLimit:          10,
				defaultCreatedAtOrder: "asc",
			},
			want: &PaginationURLQuery{
				Offset:         0,
				Limit:          10,
				CreatedAtOrder: "asc",
			},
			wantErr: false,
		},
		{
			name: "use the 0 if the passing offset is negative",
			args: args{
				rawURL:                "https://example.com?offset=-10&limit=10&created_at_order=asc",
				defaultLimit:          10,
				defaultCreatedAtOrder: "asc",
			},
			want: &PaginationURLQuery{
				Offset:         0,
				Limit:          10,
				CreatedAtOrder: "asc",
			},
			wantErr: false,
		},
		{
			name: "use the default limit if the passing is too big",
			args: args{
				rawURL:                "https://example.com?offset=0&limit=100&created_at_order=asc",
				defaultLimit:          10,
				defaultCreatedAtOrder: "asc",
			},
			want: &PaginationURLQuery{
				Offset:         0,
				Limit:          10,
				CreatedAtOrder: "asc",
			},
			wantErr: false,
		},
		{
			name: "use the default created_at_order if the passing is empty",
			args: args{
				rawURL:                "https://example.com?offset=0&limit=10",
				defaultLimit:          10,
				defaultCreatedAtOrder: "desc",
			},
			want: &PaginationURLQuery{
				Offset:         0,
				Limit:          10,
				CreatedAtOrder: "desc",
			},
			wantErr: false,
		},
		{
			name: "returns non-nil err when passing invalid offset",
			args: args{
				rawURL:                "https://example.com?offset=x&limit=10&created_at_order=asc",
				defaultLimit:          10,
				defaultCreatedAtOrder: "asc",
			},
			wantErr: true,
		},
		{
			name: "returns non-nil err when passing invalid limit",
			args: args{
				rawURL:                "https://example.com?offset=0&limit=x&created_at_order=asc",
				defaultLimit:          10,
				defaultCreatedAtOrder: "asc",
			},
			wantErr: true,
		},
		{
			name: "returns non-nil err when passing invalid created_at_order",
			args: args{
				rawURL:                "https://example.com?offset=0&limit=10&created_at_order=x",
				defaultLimit:          10,
				defaultCreatedAtOrder: "asc",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			u, err := url.Parse(tt.args.rawURL)
			if err != nil {
				t.Fatalf("invalid rawURL: %v", err)
			}

			query, getErr := parsePaginationURLQuery(u, tt.args.defaultLimit, tt.args.defaultCreatedAtOrder)
			if (getErr != nil) != tt.wantErr {
				t.Errorf("returned unexpected error: got %v, want %v", getErr, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(query, tt.want) {
				t.Errorf("returned unexpected query: got %v, want %v", query, tt.want)
			}
		})
	}
}
