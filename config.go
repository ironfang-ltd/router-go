package router

import "net/http"

type Option func(*Config)

type Config struct {
	NotFoundHandler         http.HandlerFunc
	MethodNotAllowedHandler http.HandlerFunc
}

func WithNotFoundHandler(handler http.HandlerFunc) Option {
	return func(c *Config) {
		c.NotFoundHandler = handler
	}
}

func WithMethodNotAllowedHandler(handler http.HandlerFunc) Option {
	return func(c *Config) {
		c.MethodNotAllowedHandler = handler
	}
}
