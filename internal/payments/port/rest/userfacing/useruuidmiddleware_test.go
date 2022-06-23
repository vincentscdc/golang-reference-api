package userfacing

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func Test_UserUUID(t *testing.T) {
	t.Parallel()

	type args struct {
		userUUID string
		handler  http.HandlerFunc
	}

	echoHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(getUserUUID(r.Context()).String())); err != nil {
			t.Fatalf("unexpected write into response error: %v", err)
		}
	})

	tests := []struct {
		name                   string
		args                   args
		expectedHTTPStatusCode int
		expectedBody           string
		expectedLogMsg         string
	}{
		{
			name: "happy path",
			args: args{
				userUUID: "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
				handler:  echoHandler,
			},
			expectedHTTPStatusCode: http.StatusOK,
			expectedBody:           "b7202eb0-5bf0-475d-8ee2-d3d2c168a5d5",
		},
		{
			name: "response 401 when header is not set",
			args: args{
				handler: echoHandler,
			},
			expectedHTTPStatusCode: http.StatusUnauthorized,
		},
		{
			name: "response 401 and log error when user id format is invalid",
			args: args{
				userUUID: "xxx",
				handler:  echoHandler,
			},
			expectedHTTPStatusCode: http.StatusUnauthorized,
			expectedLogMsg:         `{"level":"error","error":"invalid UUID length: 3","UserID":"xxx","message":"invalid user id"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var logMsg bytes.Buffer
			log := zerolog.New(&logMsg)

			r := chi.NewRouter()
			r.Use(UserUUID(&log))
			r.Get("/", tt.args.handler)

			srv := httptest.NewServer(r)
			defer srv.Close()

			req := httptest.NewRequest("GET", "/", nil)
			setRequestHeaderUserID(req, tt.args.userUUID)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedHTTPStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedHTTPStatusCode)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}

			//if strings.TrimSpace(logMsg.String()) != tt.expectedLogMsg {
			//	t.Errorf("handler log unexpected log message: got %v want %v", logMsg.String(), tt.expectedLogMsg)
			//}
		})
	}
}
