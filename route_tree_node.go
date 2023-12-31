package router

import (
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
	httpMethodAny
	httpMethodCount
)

type routeTreeNode struct {
	segment    string
	parent     *routeTreeNode
	children   []*routeTreeNode
	middleware []Middleware
	handlers   []http.HandlerFunc
	param      bool
	catchAll   bool
}

func newRouteTreeNode() *routeTreeNode {
	return &routeTreeNode{
		segment:    "",
		parent:     nil,
		children:   nil,
		middleware: nil,
		handlers:   nil,
		param:      false,
		catchAll:   false,
	}
}

func nodePriority(node *routeTreeNode) int {

	if node.catchAll {
		return 3
	}

	if node.param {
		return 2
	}

	// static
	return 1
}

func (r *routeTreeNode) GetOrCreateNode(path string) *routeTreeNode {

	node := r
	high := 0

	for {
		if len(path) == 0 {
			break
		}

		high = strings.IndexByte(path, PathSep)
		if high == -1 {
			high = len(path)
		}

		segment := path[:high]

		if segment == "" {
			node = r
			high++
			if high >= len(path) {
				break
			}
			path = path[high:]
			continue
		}

		found := false

		for _, child := range node.children {
			if child.segment == segment {
				// TODO: Check for conflicting param/catchAll
				node = child
				found = true
				break
			}
		}

		if !found {
			newNode := newRouteTreeNode()
			newNode.segment = segment
			newNode.parent = node
			newNode.param = segment != "" && segment[0] == ':'
			newNode.catchAll = segment != "" && segment[len(segment)-1] == '*'

			node.children = append(node.children, newNode)

			sort.Slice(node.children, func(i, j int) bool {
				// Sort Order: segment(static) > param > catchAll
				return nodePriority(node.children[i]) < nodePriority(node.children[j])
			})

			node = newNode
		}

		if segment == "*" {
			break
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
			} else if child.catchAll {
				return child, params
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

	if r.handlers == nil {
		return nil
	}

	if r.handlers[httpMethodAny] != nil {
		return r.handlers[httpMethodAny]
	}

	return r.handlers[methodToUint8(method)]
}

func (r *routeTreeNode) Use(m ...Middleware) {

	if r.middleware == nil {
		r.middleware = make([]Middleware, 0)
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
	case "*":
		return httpMethodAny
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
	case httpMethodAny:
		return "*"
	default:
		panic("unhandled default case")
	}
}
