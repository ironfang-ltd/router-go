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

	var result error

	m := Recover(WithHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusInternalServerError)

		result = err
	}))

	m(w, req, func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("something went wrong"))
	})

	if result == nil {
		t.Fatalf("expected error, got nil")
	}
}
