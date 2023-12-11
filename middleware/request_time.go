package middleware

import (
	"net/http"
	"strconv"
	"time"
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

func RequestTime(opts ...RequestTimeOption) Handler {

	config := &RequestTimeConfig{
		HeaderName: "X-Request-Time-Ms",
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(w http.ResponseWriter, r *http.Request, next Next) error {

		start := time.Now()
		err := next(w, r)
		taken := time.Now().Sub(start)

		w.Header().Set(config.HeaderName, strconv.FormatInt(taken.Milliseconds(), 10))

		return err
	}
}
