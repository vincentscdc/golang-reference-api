package userfacing

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/monacohq/golang-common/transport/http/handlerwrap/v2"
	"github.com/rs/zerolog"
)

func TestOKStyleWrapper(t *testing.T) {
	t.Parallel()

	type args struct {
		dataName string
		handler  handlerwrap.TypedHandler
	}

	tests := []struct {
		name         string
		args         args
		logMsg       string
		responseBody string
	}{
		{
			name: "happy path: nil error",
			args: args{
				dataName: "foo",
				handler: func(r *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
					return &handlerwrap.Response{
						Body:       map[string]any{"bar": "bar"},
						StatusCode: http.StatusCreated,
					}, nil
				},
			},
			logMsg:       "",
			responseBody: `{"foo":{"bar":"bar"},"ok":true}`,
		},
		{
			name: "returns nil error and log error when the handler returns a non-nil error",
			args: args{
				dataName: "foo",
				handler: func(r *http.Request) (*handlerwrap.Response, *handlerwrap.ErrorResponse) {
					return &handlerwrap.Response{
							Body:       map[string]any{"bar": "bar"},
							StatusCode: http.StatusBadRequest,
						}, &handlerwrap.ErrorResponse{
							Error:      errors.New("bad request"),
							StatusCode: http.StatusBadRequest,
							ErrorCode:  "bad_request",
							ErrorMsg:   "bad request",
						}
				},
			},
			logMsg:       `{"level":"error","error":"bad request","ErrorCode":"bad_request","HTTPStatusCode":400,"message":"bad request"}`,
			responseBody: `{"error":"bad_request","error_message":"bad request","ok":false}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var logBuffer bytes.Buffer
			log := zerolog.New(&logBuffer)

			req := httptest.NewRequest("GET", "/", nil)
			resp, err := OKStyleWrapper(&log, tt.args.dataName, tt.args.handler)(req)
			if err != nil {
				t.Fatalf("handler returned a non-nil error, got %v want nil", err)
			}
			if msg := strings.TrimSpace(logBuffer.String()); msg != tt.logMsg {
				t.Errorf("handler loged a wrong message, got %v want %v", msg, tt.logMsg)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("handler returned wrong status code, got %v want %v",
					http.StatusOK, resp.StatusCode)
			}
			respJSON, _ := json.Marshal(resp.Body)
			if string(respJSON) != tt.responseBody {
				t.Errorf("handler returned wrong body content, got %v want %v",
					string(respJSON), tt.responseBody)
			}
		})
	}
}
