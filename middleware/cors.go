package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ironfang-ltd/router-go"
)

type CorsOption func(*CorsConfig)

type CorsConfig struct {
	Origins          []string
	Methods          []string
	Headers          []string
	MaxAge           int
	AllowCredentials bool
}

func WithOrigins(origins ...string) CorsOption {
	return func(c *CorsConfig) {
		c.Origins = origins
	}
}

func WithMethods(methods ...string) CorsOption {
	return func(c *CorsConfig) {
		c.Methods = methods
	}
}

func WithHeaders(headers ...string) CorsOption {
	return func(c *CorsConfig) {
		c.Headers = headers
	}
}

func WithCredentials() CorsOption {
	return func(c *CorsConfig) {
		c.AllowCredentials = true
	}
}

func Cors(opts ...CorsOption) router.Middleware {

	config := &CorsConfig{
		Origins: []string{"*"},
		Methods: []string{"OPTIONS", "PUT", "PATCH", "DELETE"},
		Headers: []string{"Content-Type", "Authorization"},
		MaxAge:  3600, // 1 Hour in Seconds
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		w.Header().Add("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

		origin := r.Header.Get("Origin")

		if origin != "" {

			// Loop through the list of trusted origins, checking to see if the request
			// origin exactly matches one of them.
			for i := range config.Origins {
				if origin == config.Origins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if config.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

						w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Methods, ", "))
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Headers, ", "))
						w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))

						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next(w, r)
	}
}
