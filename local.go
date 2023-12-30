package router

import (
	"context"
	"net/http"
)

type RequestLocals map[string]interface{}

func (l RequestLocals) Get(key string) interface{} {
	return l[key]
}

func (l RequestLocals) Set(key string, value interface{}) {
	l[key] = value
}

func SetLocal(r *http.Request, key string, value interface{}) *http.Request {

	ctx := r.Context()

	if ctx == nil {
		ctx = context.Background()
	}

	locals := ctx.Value(contextKeyLocal)

	if locals == nil {
		locals = RequestLocals{}
	}

	locals.(RequestLocals).Set(key, value)

	ctx = context.WithValue(ctx, contextKeyLocal, locals)

	return r.WithContext(ctx)
}

func GetLocal(r *http.Request, key string) interface{} {

	ctx := r.Context()

	if ctx != nil {
		locals := ctx.Value(contextKeyLocal)

		if locals != nil {
			return locals.(RequestLocals).Get(key)
		}
	}

	return nil
}
