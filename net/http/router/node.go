package router

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.001
// @date    2018-10-02

import (
	hdr "github.com/lj-team/go-generic/net/http/handler"
	"github.com/lj-team/go-generic/net/http/url/params"
)

type (
	Handler = hdr.Handler
	Params  = params.Params
)

type node struct {
	Name   string
	Childs map[string]*node
	F      Handler
}
