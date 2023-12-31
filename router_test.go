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

func TestRouter_Static(t *testing.T) {

	req, _ := http.NewRequest("GET", "/test.txt", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Static("/", "./testdata")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("response content type is not text/plain")
	}

	if w.Body.String() != "hello world" {
		t.Error("response body is not 'hello world'")
	}
}

func TestRouter_StaticWithPath(t *testing.T) {

	req, _ := http.NewRequest("GET", "/files/test.txt", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Static("/files", "./testdata")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("response content type is not text/plain")
	}

	if w.Body.String() != "hello world" {
		t.Error("response body is not 'hello world'")
	}
}

func TestRouter_NodeOrder(t *testing.T) {

	r := New()

	r.Static("/files", "./testdata")

	r.Get("/:param", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {})

	routes := r.GetRoutes()

	if routes[0].Path != "/test" {
		t.Error("first route is not '/test', but ", routes[0].Path)
	}

	if routes[1].Path != "/:param" {
		t.Error("second route is not '/:param', but ", routes[1].Path)
	}

	if routes[2].Path != "/files*" {
		t.Error("third route is not '/files*', but ", routes[2].Path)
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
