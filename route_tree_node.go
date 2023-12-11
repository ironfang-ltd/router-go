package router

import (
	"github.com/ironfang-ltd/router-go/middleware"
	"net/http"
	"sort"
	"strings"
)

const (
	PathSep = '/'
)

const (
	httpMethodGet uint8 = iota
	httpMethodHead
	httpMethodPost
	httpMethodPut
	httpMethodPatch
	httpMethodDelete
	httpMethodConnect
	httpMethodOptions
	httpMethodTrace
	httpMethodCount
)

type routeTreeNode struct {
	segment    string
	parent     *routeTreeNode
	children   []*routeTreeNode
	middleware []middleware.Handler
	handlers   []http.HandlerFunc
	param      bool
}

func newRouteTreeNode() *routeTreeNode {
	return &routeTreeNode{
		segment:    "",
		parent:     nil,
		children:   nil,
		middleware: nil,
		handlers:   nil,
		param:      false,
	}
}

func (r *routeTreeNode) GetOrCreateNode(path string) *routeTreeNode {

	if path == "" {
		return r
	}

	if path == "/" {
		return r
	}

	node := r
	high := 0

	if path[0] == PathSep {
		path = path[1:]
	}

	for {
		if len(path) == 0 {
			break
		}

		high = strings.IndexByte(path, PathSep)
		if high == -1 {
			high = len(path)
		}

		segment := path[:high]
		found := false

		for _, child := range node.children {
			if child.segment == segment {
				node = child
				found = true
				break
			}
		}

		if !found {
			newNode := newRouteTreeNode()
			newNode.segment = segment
			newNode.parent = node
			newNode.param = segment[0] == ':'

			node.children = append(node.children, newNode)

			sort.SliceStable(node.children, func(i, j int) bool {
				return !node.param
			})

			node = newNode
		}

		if high >= len(path) {
			break
		}

		high++
		path = path[high:]
	}

	return node
}

func (r *routeTreeNode) Find(path string) (*routeTreeNode, *routeParams) {

	if path == "" {
		return nil, nil
	}

	if path == "/" {
		return r, nil
	}

	node := r
	high := 0

	if path[0] == PathSep {
		path = path[1:]
	}

	params := &routeParams{}

	for {
		if len(path) == 0 {
			break
		}

		high = strings.IndexByte(path, PathSep)
		if high == -1 {
			high = len(path)
		}

		segment := path[:high]
		high++
		found := false

		for _, child := range node.children {

			if child.param {

				params.set(child.segment[1:], segment)

				if high >= len(path) {
					return child, params
				}

				node = child
				found = true
				path = path[high:]
				break
			} else if child.segment == segment {
				if high >= len(path) {
					return child, params
				}

				node = child
				found = true
				path = path[high:]
				break
			}
		}

		if !found {
			return nil, nil
		}
	}

	return node, params
}

func (r *routeTreeNode) SetHandler(method string, handler http.HandlerFunc) {
	if r.handlers == nil {
		r.handlers = make([]http.HandlerFunc, httpMethodCount)
	}

	r.handlers[methodToUint8(method)] = handler
}

func (r *routeTreeNode) GetHandler(method string) http.HandlerFunc {
	return r.handlers[methodToUint8(method)]
}

func (r *routeTreeNode) Use(m ...middleware.Handler) {

	if r.middleware == nil {
		r.middleware = make([]middleware.Handler, 0)
	}

	r.middleware = append(r.middleware, m...)
}

func (r *routeTreeNode) getPath() string {

	if r.parent == nil {
		return r.segment
	}

	return r.parent.getPath() + "/" + r.segment
}

func methodToUint8(method string) uint8 {

	switch method {
	case http.MethodGet:
		return httpMethodGet
	case http.MethodHead:
		return httpMethodHead
	case http.MethodPost:
		return httpMethodPost
	case http.MethodPut:
		return httpMethodPut
	case http.MethodPatch:
		return httpMethodPatch
	case http.MethodDelete:
		return httpMethodDelete
	case http.MethodConnect:
		return httpMethodConnect
	case http.MethodOptions:
		return httpMethodOptions
	case http.MethodTrace:
		return httpMethodTrace
	}

	return httpMethodGet
}

func uint8ToMethod(method uint8) string {

	switch method {
	case httpMethodGet:
		return http.MethodGet
	case httpMethodHead:
		return http.MethodHead
	case httpMethodPost:
		return http.MethodPost
	case httpMethodPut:
		return http.MethodPut
	case httpMethodPatch:
		return http.MethodPatch
	case httpMethodDelete:
		return http.MethodDelete
	case httpMethodConnect:
		return http.MethodConnect
	case httpMethodOptions:
		return http.MethodOptions
	case httpMethodTrace:
		return http.MethodTrace
	default:
		panic("unhandled default case")
	}

	return http.MethodGet
}
