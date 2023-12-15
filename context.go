package router

type contextKey string

var (
	contextKeyRoute = contextKey("route")
	contextKeyLocal = contextKey("local")
)
