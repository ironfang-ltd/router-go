package router

import "net/http"

type contextKey string

var (
	contextKeyRoute = contextKey("route")
)

func RouteParam(r *http.Request, name string) string {

	ctx := r.Context()

	if ctx != nil {
		route := ctx.Value(contextKeyRoute)

		if route != nil {
			p := route.(*routeParams)

			return p.get(name)
		}
	}

	return ""
}
