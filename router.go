package router

import (
	"context"
	"fmt"
	"github.com/ironfang-ltd/router-go/middleware"
	"net/http"
)

const (
	ErrPathMustStartWithSlash  = "path must start with '/'"
	ErrPathMustNotEndWithSlash = "path must not end with '/'"
)

type Route interface {
	Use(middleware ...middleware.Handler)
}

type Router interface {
	Get(path string, handler http.HandlerFunc) Route
	Post(path string, handler http.HandlerFunc) Route
	Put(path string, handler http.HandlerFunc) Route
	Patch(path string, handler http.HandlerFunc) Route
	Delete(path string, handler http.HandlerFunc) Route
	Group(prefix string) Router
	Use(middleware ...middleware.Handler)
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
}

func New() Router {
	r := router{
		parent: nil,
		prefix: "",
		root:   newRouteTreeNode(),
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

func (r *router) Use(m ...middleware.Handler) {
	r.root.Use(m...)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

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

	err := r.handleMiddleware(node, w, req, handler)
	if err != nil {
		r.internalError(w, req, err)
		return
	}
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

		for _, child := range node.children {
			q = append(q, child)
		}
	}

	return routes
}

func (r *router) handleMiddleware(n *routeTreeNode, w http.ResponseWriter, req *http.Request, final http.HandlerFunc) error {

	if n.parent != nil {
		err := r.handleMiddleware(n.parent, w, req, final)
		if err != nil {
			return err
		}

		return nil
	}

	if n.middleware == nil {
		final(w, req)

		return nil
	}

	mc := middlewareContext{
		current:    0,
		middleware: n.middleware,
		final:      final,
	}

	err := mc.Next(w, req)
	if err != nil {
		return err
	}

	return nil
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

func (r *router) methodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (r *router) notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (r *router) internalError(w http.ResponseWriter, _ *http.Request, err error) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}
