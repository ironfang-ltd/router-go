package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Timeout is a middleware that cancels the ctx after a given timeout
// and returns HTTP Status 504 (Gateway Timeout).
func Timeout(timeout time.Duration) Handler {

	return func(w http.ResponseWriter, r *http.Request, next Next) error {

		ctx, cancel := context.WithTimeout(r.Context(), timeout)

		defer func() {

			// Calling cancel() on a canceled context is a noop.
			cancel()

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		}()

		return next(w, r.WithContext(ctx))
	}
}
