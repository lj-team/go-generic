package handler

import (
	"net/http"

	"github.com/lj-team/go-generic/net/http/url/params"
)

type (
	Params = params.Params
)

type Handler func(rw http.ResponseWriter, r *http.Request, params Params)
