package router

import (
	"net/http"

	"github.com/lj-team/go-generic/net/http/errors"
)

var (
	handler  *Router
	initTime int64
)

func SetDefault(rt *Router) {
	handler = rt
}

func GetDefault() *Router {
	return handler
}

func Register(method string, path string, fn Handler) {
	if handler != nil {
		handler.Register(method, path, fn)
	}
}

func Redirect(from, to string, code int) {
	if handler != nil {
		handler.Redirect(from, to, code)
	}
}

func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if handler != nil {
		handler.ServeHTTP(rw, req)
	} else {
		errors.Send(rw, 404)
	}
}
