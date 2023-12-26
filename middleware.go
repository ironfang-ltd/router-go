package router

import (
	"net/http"
)

type Middleware func(http.ResponseWriter, *http.Request, http.HandlerFunc)

type middlewareContext struct {
	current    int
	middleware []Middleware
	final      http.HandlerFunc
}

func (mc *middlewareContext) Next(w http.ResponseWriter, req *http.Request) {

	if mc.current >= len(mc.middleware) {

		if mc.final != nil {
			mc.final(w, req)
		}

		return
	}

	c := mc.current
	mc.current++

	mc.middleware[c](w, req, mc.Next)
}
