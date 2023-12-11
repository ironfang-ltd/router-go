package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_GetWithParam(t *testing.T) {

	req, _ := http.NewRequest("GET", "/test-value", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Get("/:param", func(w http.ResponseWriter, r *http.Request) {

		param := RouteParam(r, "param")

		_, _ = w.Write([]byte(param))
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Body.String() != "test-value" {
		t.Error("response body is not 'test-value'")
	}
}

func BenchmarkGet(b *testing.B) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGetWithParam(b *testing.B) {

	req, _ := http.NewRequest("GET", "/test-value", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Get("/:name", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(RouteParam(r, "name")))
	})

	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}
