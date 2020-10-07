package router

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lj-team/go-generic/log"
	"github.com/lj-team/go-generic/net/http/errors"
)

func init() {

	initTime = time.Now().Unix()

	New(true)
}

type Router struct {
	methods           map[string]*node
	redirects         map[string]string
	NotFoundFunc      http.HandlerFunc
	BadRequestFunc    http.HandlerFunc
	InternalErrorFunc http.HandlerFunc
	OptionsFunc       http.HandlerFunc
	MakeUid           bool
}

func New(isDefault bool) *Router {
	rt := &Router{
		methods:           map[string]*node{},
		redirects:         map[string]string{},
		NotFoundFunc:      func(rw http.ResponseWriter, req *http.Request) { errors.Send(rw, 404) },
		BadRequestFunc:    func(rw http.ResponseWriter, req *http.Request) { errors.Send(rw, 400) },
		InternalErrorFunc: func(rw http.ResponseWriter, req *http.Request) { errors.Send(rw, 500) },
		OptionsFunc:       nil,
		MakeUid:           true,
	}

	if isDefault {
		handler = rt
	}

	rt.Register("GET", "/alive", aliveHandler)
	rt.Register("POST", "/alive", aliveHandler)

	return rt
}

func (r *Router) Register(method string, path string, fn Handler) {

	root, has := r.methods[method]
	if !has {
		root = &node{Childs: make(map[string]*node)}
		r.methods[method] = root
	}

	list := path2list(path, true)
	if list == nil {
		return
	}

	for _, item := range list {

		if item[0] == ':' {
			n, h := root.Childs[""]
			if !h {
				name := item
				if len(name) > 1 {
					name = item[1:]
				} else {
					name = ""
				}
				n = &node{Name: name, Childs: make(map[string]*node)}
			}
			root.Childs[""] = n
			root = n
			continue
		}

		if item == "*" {
			_, h := root.Childs["*"]
			if !h {
				root.Childs["*"] = &node{Name: "*", F: fn}
			}
			return
		}

		n, h := root.Childs[item]
		if !h {
			n = &node{Name: "", Childs: make(map[string]*node)}
		}
		root.Childs[item] = n
		root = n

	}

	root.F = fn
}

func (r *Router) Redirect(from, to string, code int) {
	r.redirects[from] = to
	r.redirects[from+"/"] = to
}

func searchHandler(n *node, paths []string, params []string) (Handler, []string) {

	switch len(paths) {

	case 0:
		return nil, nil

	case 1:

		key := paths[0]

		if nn, h := n.Childs[key]; h {
			return nn.F, params
		}

		if nn, h := n.Childs[""]; h {
			params = append(params, nn.Name, key)
			return nn.F, params
		}

		if nn, h := n.Childs["*"]; h {
			params = append(params, "*", "/"+key)
			return nn.F, params
		}

	default:

		key := paths[0]

		if nn, h := n.Childs[key]; h {
			f, p := searchHandler(nn, paths[1:], params)
			if f != nil {
				return f, p
			}
		}

		if nn, h := n.Childs[""]; h {
			params = append(params, nn.Name, key)
			f, p := searchHandler(nn, paths[1:], params)
			if f != nil {
				return f, p
			}
			params = params[:len(params)-2]
		}

		if nn, h := n.Childs["*"]; h {
			params = append(params, "*", "/"+strings.Join(paths, "/"))
			return nn.F, params
		}
	}

	return nil, nil
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		if dest, has := r.redirects[req.URL.Path]; has {
			log.Trace("GET " + req.URL.Path + " -> " + dest)
			http.Redirect(rw, req, dest, 302)
			return
		}
	}

	writeLog(req.Method, req.RequestURI)

	defer func() {

		if re := recover(); re != nil {
			log.Error(fmt.Sprint(re))
			log.Stack("error")
			r.InternalErrorFunc(rw, req)
		}

	}()

	if req.Method == http.MethodOptions && r.OptionsFunc != nil {
		r.OptionsFunc(rw, req)
		return
	}

	root, has := r.methods[req.Method]
	if !has {
		r.NotFoundFunc(rw, req)
		return
	}

	fn, params := searchHandler(root, path2list(req.URL.Path, false), []string{})

	if fn != nil {

		prm := map[string]string{}

		for i := 0; i+1 < len(params); i += 2 {
			prm[params[i]] = params[i+1]
		}

		fn(rw, req, Params(prm))
		return
	}

	r.NotFoundFunc(rw, req)
}

func path2list(path string, regmod bool) []string {
	rh := strings.NewReader(path)
	builder := strings.Builder{}
	mode := 0
	result := make([]string, 0, 8)

	for {
		run, _, err := rh.ReadRune()
		if err != nil {
			break
		}

		switch mode {
		case 0:

			if run != '/' {
				return nil
			}

			mode = 1
			result = append(result, "/")

		case 1:

			if run != '/' {
				builder.WriteRune(run)
				mode = 2
			}

		case 2:

			if run == '/' {
				str := builder.String()
				builder.Reset()

				if str == "*" {
					result = append(result, "*")
					if regmod {
						mode = 3
					} else {
						return nil
					}
				} else {

					uri, e := url.PathUnescape(str)
					if e != nil {
						return nil
					}

					result = append(result, uri)
					mode = 1
				}
			} else {
				builder.WriteRune(run)
			}

		case 3:

			if run != '/' {
				return nil
			}
		}
	}

	switch mode {

	case 0:
		return nil

	case 2:
		str := builder.String()
		builder.Reset()

		uri, e := url.PathUnescape(str)
		if e != nil {
			return nil
		}

		result = append(result, uri)
	}

	return result
}

func writeLog(method string, url string) {
	log.Trace(fmt.Sprintf("%s %s", method, url))
}

func aliveHandler(rw http.ResponseWriter, req *http.Request, prms Params) {

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	fmt.Fprintf(rw, "{\"alive\":%d}", time.Now().Unix()-initTime)
}
