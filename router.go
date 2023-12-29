package router

import (
	"context"
	"net/http"
)

const (
	ErrPathMustStartWithSlash  = "path must start with '/'"
	ErrPathMustNotEndWithSlash = "path must not end with '/'"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type Route interface {
	Use(middleware ...Middleware)
}

type Router interface {
	Get(path string, handler http.HandlerFunc) Route
	Post(path string, handler http.HandlerFunc) Route
	Put(path string, handler http.HandlerFunc) Route
	Patch(path string, handler http.HandlerFunc) Route
	Delete(path string, handler http.HandlerFunc) Route
	Group(prefix string) Router
	Use(middleware ...Middleware)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetRoutes() []RouteDescriptor
}

type RouteDescriptor struct {
	Method string
	Path   string
}

type router struct {
	parent *router
	prefix string
	root   *routeTreeNode
	config *RouterConfig
}

func New(opts ...RouterOption) Router {

	config := &RouterConfig{
		NotFoundHandler:         nil,
		MethodNotAllowedHandler: nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	r := router{
		parent: nil,
		prefix: "",
		root:   newRouteTreeNode(),
		config: config,
	}

	return &r
}

func (r *router) Get(path string, handler http.HandlerFunc) Route {
	return r.mapMethod(http.MethodGet, path, handler)
}

func (r *router) Post(path string, handler http.HandlerFunc) Route {
	return r.mapMethod(http.MethodPost, path, handler)
}

func (r *router) Put(path string, handler http.HandlerFunc) Route {
	return r.mapMethod(http.MethodPut, path, handler)
}

func (r *router) Patch(path string, handler http.HandlerFunc) Route {
	return r.mapMethod(http.MethodPatch, path, handler)
}

func (r *router) Delete(path string, handler http.HandlerFunc) Route {
	return r.mapMethod(http.MethodDelete, path, handler)
}

func (r *router) Group(prefix string) Router {
	group := router{
		parent: r,
		prefix: prefix,
		root:   nil,
	}

	return &group
}

func (r *router) Use(m ...Middleware) {
	r.root.Use(m...)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	r.handleMiddleware(r.root, w, req, func(w http.ResponseWriter, req *http.Request) {

		node, params := r.root.Find(req.URL.Path)
		if node == nil {
			r.notFound(w, req)
			return
		}

		handler := node.GetHandler(req.Method)
		if handler == nil {
			r.methodNotAllowed(w, req)
			return
		}

		if params != nil {
			ctx := context.WithValue(req.Context(), contextKeyRoute, params)
			req = req.WithContext(ctx)
		}

		r.handleMiddleware(node, w, req, handler)
	})
}

func (r *router) GetRoutes() []RouteDescriptor {

	var routes []RouteDescriptor

	q := []*routeTreeNode{r.root}

	for {
		if len(q) == 0 {
			break
		}

		node := q[0]
		q = q[1:]

		if node == nil {
			break
		}

		for i, handler := range node.handlers {
			if handler != nil {

				p := node.getPath()
				if len(p) == 0 {
					p = "/"
				}

				routes = append(routes, RouteDescriptor{
					Method: uint8ToMethod(uint8(i)),
					Path:   p,
				})
			}
		}

		q = append(q, node.children...)
	}

	return routes
}

func (r *router) handleMiddleware(n *routeTreeNode, w http.ResponseWriter, req *http.Request, final http.HandlerFunc) {

	if n.parent != nil && n.parent != r.root {
		r.handleMiddleware(n.parent, w, req, final)

		return
	}

	if n.middleware == nil {
		final(w, req)

		return
	}

	mc := middlewareContext{
		current:    0,
		middleware: n.middleware,
		final:      final,
	}

	mc.Next(w, req)
}

func (r *router) mapMethod(method, path string, handler http.HandlerFunc) *routeTreeNode {

	if r.parent != nil {
		return r.parent.mapMethod(method, r.prefix+path, handler)
	}

	if len(path) == 0 || path[0] != PathSep {
		panic(ErrPathMustStartWithSlash)
	}

	if len(path) > 1 && path[len(path)-1] == PathSep {
		panic(ErrPathMustNotEndWithSlash)
	}

	node := r.root.GetOrCreateNode(path)
	node.SetHandler(method, handler)
	return node
}

func (r *router) methodNotAllowed(w http.ResponseWriter, req *http.Request) {

	if r.config.MethodNotAllowedHandler != nil {
		r.config.MethodNotAllowedHandler(w, req)

		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (r *router) notFound(w http.ResponseWriter, req *http.Request) {

	if r.config.NotFoundHandler != nil {
		r.config.NotFoundHandler(w, req)

		return
	}

	w.WriteHeader(http.StatusNotFound)
}
