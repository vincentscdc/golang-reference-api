package userfacing

import (
	"flag"
	"net/http"
	"os"
	"testing"

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

func setRequestHeaderUserID(r *http.Request, uuid string) {
	r.Header.Set(HTTPHeaderKeyUserUUID, uuid)
}
