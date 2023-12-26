package middleware

import (
	"fmt"
	"net/http"

	"github.com/ironfang-ltd/router-go"
)

type RecoverOption func(*RecoverConfig)

type RecoverConfig struct {
	ErrorHandler router.ErrorHandler
}

func WithHandler(handler router.ErrorHandler) RecoverOption {
	return func(c *RecoverConfig) {
		c.ErrorHandler = handler
	}
}

func Recover(opts ...RecoverOption) router.Middleware {

	config := &RecoverConfig{
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		defer func() {
			if rerr := recover(); rerr != nil {

				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")

				if e, ok := rerr.(error); ok {
					config.ErrorHandler(w, r, e)
				} else {
					config.ErrorHandler(w, r, fmt.Errorf("%v", rerr))
				}
			}
		}()

		next(w, r)
	}
}
