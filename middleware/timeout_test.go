package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	m := Timeout(1 * time.Second)

	m(w, req, func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
			t.Fatal("timeout middleware failed to timeout the request")
		}

	})

	if w.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status code %d, got %d", http.StatusGatewayTimeout, w.Code)
	}
}
