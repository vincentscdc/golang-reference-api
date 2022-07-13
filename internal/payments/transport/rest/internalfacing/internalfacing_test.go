package internalfacing

import (
	"context"
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func setURLParams(req *http.Request, params map[string]string) {
	ctx := chi.NewRouteContext()

	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}

	*req = *req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (int, error) {
	return rf(p)
}
