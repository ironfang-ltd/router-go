package router

import "net/http"

type RouterOption func(*RouterConfig)

type RouterConfig struct {
	NotFoundHandler         http.HandlerFunc
	MethodNotAllowedHandler http.HandlerFunc
}

func WithNotFoundHandler(handler http.HandlerFunc) RouterOption {
	return func(c *RouterConfig) {
		c.NotFoundHandler = handler
	}
}

func WithMethodNotAllowedHandler(handler http.HandlerFunc) RouterOption {
	return func(c *RouterConfig) {
		c.MethodNotAllowedHandler = handler
	}
}
