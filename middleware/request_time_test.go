package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestTime(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	m := RequestTime()

	m(w, req, func(w http.ResponseWriter, r *http.Request) {

	})

	if w.Header().Get("X-Request-Time-Ms") == "" {
		t.Fatal("expected header X-Request-Time-Ms to be set")
	}
}

func TestRequestTimeWithCustomHeader(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	m := RequestTime(WithHeaderName("X-Request-Time"))

	m(w, req, func(w http.ResponseWriter, r *http.Request) {
	})

	if w.Header().Get("X-Request-Time") == "" {
		t.Fatal("expected custom header X-Request-Time to be set")
	}
}
