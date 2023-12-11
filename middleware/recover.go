package middleware

import (
	"fmt"
	"net/http"
)

func Recover() Handler {

	return func(w http.ResponseWriter, r *http.Request, next Next) (err error) {

		defer func() {
			if rerr := recover(); rerr != nil {

				if e, ok := rerr.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", rerr)
				}
			}
		}()

		return next(w, r)
	}
}
