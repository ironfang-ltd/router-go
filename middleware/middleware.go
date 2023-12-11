package middleware

import "net/http"

type Next func(http.ResponseWriter, *http.Request) error
type Handler func(http.ResponseWriter, *http.Request, Next) error
