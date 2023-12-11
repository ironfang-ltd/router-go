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

	err := m(w, req, func(w http.ResponseWriter, r *http.Request) error {

		ctx := r.Context()

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(2 * time.Second):
			t.Fatal("timeout middleware failed to timeout the request")
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if w.Code != http.StatusGatewayTimeout {
		t.Fatalf("expected status code %d, got %d", http.StatusGatewayTimeout, w.Code)
	}
}
