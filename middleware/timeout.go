package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ironfang-ltd/router-go"
)

// Timeout is a middleware that cancels the ctx after a given timeout
// and returns HTTP Status 504 (Gateway Timeout).
func Timeout(timeout time.Duration) router.Middleware {

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		ctx, cancel := context.WithTimeout(r.Context(), timeout)

		defer func() {

			// Calling cancel() on a canceled context is a noop.
			cancel()

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		}()

		next(w, r.WithContext(ctx))
	}
}
