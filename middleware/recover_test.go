package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecover(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	m := Recover()

	err := m(w, req, func(w http.ResponseWriter, r *http.Request) error {
		panic("something went wrong")
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

}

func TestRecoverWithError(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	m := Recover()

	err := m(w, req, func(w http.ResponseWriter, r *http.Request) error {
		panic(fmt.Errorf("something went wrong"))
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
