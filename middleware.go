package router

import (
	"github.com/ironfang-ltd/router-go/middleware"
	"net/http"
)

type middlewareContext struct {
	current    int
	middleware []middleware.Handler
	complete   bool
}

func (mc *middlewareContext) Next(w http.ResponseWriter, req *http.Request) error {

	if mc.current >= len(mc.middleware) {
		mc.complete = true
		return nil
	}

	c := mc.current
	mc.current++

	err := mc.middleware[c](w, req, mc.Next)
	if err != nil {
		return err
	}

	return nil
}
