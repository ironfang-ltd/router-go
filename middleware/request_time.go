package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ironfang-ltd/router-go"
)

type RequestTimeOption func(*RequestTimeConfig)

type RequestTimeConfig struct {
	HeaderName string
}

func WithHeaderName(name string) RequestTimeOption {
	return func(c *RequestTimeConfig) {
		c.HeaderName = name
	}
}

func RequestTime(opts ...RequestTimeOption) router.Middleware {

	config := &RequestTimeConfig{
		HeaderName: "X-Request-Time-Ms",
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		start := time.Now()
		next(w, r)
		taken := time.Since(start)

		w.Header().Set(config.HeaderName, strconv.FormatInt(taken.Milliseconds(), 10))
	}
}
