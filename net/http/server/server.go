package server

import (
	"net/http"

	"github.com/lj-team/go-generic/log"
	"github.com/lj-team/go-generic/net/http/handler"
	"github.com/lj-team/go-generic/net/http/router"
	"github.com/lj-team/go-generic/tmpl"
)

type (
	HandlerFunc = handler.Handler
	TmplVars    = tmpl.VarMap
)

func Start(addr string) error {

	log.InfoParams("start server on", addr)

	return http.ListenAndServe(addr, router.GetDefault())
}

func RegHandler(method string, url string, fn HandlerFunc) {
	router.GetDefault().Register(method, url, fn)
}

func VarsTT() TmplVars {
	return tmpl.Vars()
}

func RenderStringTT(ttStr string, vars TmplVars) (string, error) {

	tt := tmpl.New([]string{}, false)

	if err := tt.LoadTmplString("tt", ttStr); err != nil {
		return "", err
	}

	return tt.Render("tt", vars)
}

func WriteStringTT(rw http.ResponseWriter, ttStr string, vars TmplVars) {
	str, err := RenderStringTT(ttStr, vars)
	if err != nil {
		log.Error(err.Error())
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(str))
}
