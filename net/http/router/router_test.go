package router

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.005
// @date    2019-08-02

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPath2List(t *testing.T) {

	data := map[string][]string{
		"":                   nil,
		"/":                  []string{"/"},
		"//":                 []string{"/"},
		"///":                []string{"/"},
		"/test":              []string{"/", "test"},
		"/test//":            []string{"/", "test"},
		"/hello/world":       []string{"/", "hello", "world"},
		"/hello/world/":      []string{"/", "hello", "world"},
		"/hello/mr/:name":    []string{"/", "hello", "mr", ":name"},
		"/hello/tail/*":      []string{"/", "hello", "tail", "*"},
		"/hello/tail/*/":     []string{"/", "hello", "tail", "*"},
		"/hello/tail/*//":    []string{"/", "hello", "tail", "*"},
		"/hello/tail/*/test": nil,
	}

	for k, v := range data {
		res := path2list(k, true)

		if len(res) != len(v) || (res == nil && v != nil) || (res != nil && v == nil) {
			t.Fatal("failed " + k)
		}

		for i, item := range res {
			if item != v[i] {
				t.Fatal("failed " + k)
			}
		}
	}
}

var funcNum int = 0
var funcUid int64 = 0
var funcPid int64 = 0
var funcTail string = ""

func TestRouter(t *testing.T) {

	rt := New(false)

	rt.OptionsFunc = func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("OPTIONS"))
	}

	rt.Register("GET", "/", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/"))
	})

	rt.Register("GET", "/arg", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/arg"))
	})

	rt.Register("GET", "/*", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/*"))
	})

	rt.Register("GET", "/arg2/:param3/arg4", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/arg2/:param3/arg4|param3=" + p.GetString("param3")))
	})

	rt.Register("GET", "/arg2/arg5/arg4", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/arg2/arg5/arg4"))
	})

	rt.Register("GET", "/arg2/arg5/arg4/*", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/arg2/arg5/arg4/*|*=" + p.GetString("*")))
	})

	rt.Register("POST", "/:p1/*", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/:p1/*|*=" + p.GetString("*")))
	})

	rt.Register("POST", "/post/:p2", func(rw http.ResponseWriter, req *http.Request, p Params) {
		rw.WriteHeader(200)
		rw.Write([]byte("/post/:p2|p2=" + p.GetString("p2")))
	})

	tF := func(m string, uri string, code int, answer string) {

		rw := httptest.NewRecorder()

		req, _ := http.NewRequest(m, uri, nil)

		rt.ServeHTTP(rw, req)

		res := rw.Result()

		if res.StatusCode != code {
			t.Fatal("Invalid status code")
		}

		if res.StatusCode != 200 {
			return
		}

		body, _ := ioutil.ReadAll(res.Body)

		if string(body) != answer {
			t.Fatal("invalid answer body")
		}
	}

	tF("GET", "/", 200, "/")
	tF("GET", "/arg", 200, "/arg")
	tF("GET", "/arg/", 200, "/arg")
	tF("GET", "/arg1", 200, "/*")
	tF("GET", "/about", 200, "/*")
	tF("GET", "/arg1/arg2", 200, "/*")
	tF("GET", "/arg2/arg3/arg4", 200, "/arg2/:param3/arg4|param3=arg3")
	tF("GET", "/arg2/arg5/arg4", 200, "/arg2/arg5/arg4")
	tF("GET", "/arg2/arg5/arg4/", 200, "/arg2/arg5/arg4")
	tF("GET", "/arg2/arg5/arg4/18/13", 200, "/arg2/arg5/arg4/*|*=/18/13")
	tF("GET", "/arg2/arg3/arg4/1", 200, "/*")
	tF("GET", "/alive", 200, "{\"alive\":0}")
	tF("DELETE", "/", 404, "")
	tF("OPTIONS", "/", 200, "OPTIONS")
	tF("POST", "/t1/", 404, "")
	tF("POST", "/t1/t2", 200, "/:p1/*|*=/t2")
	tF("POST", "/post/t2", 200, "/post/:p2|p2=t2")
	tF("POST", "/post/t2/t3", 200, "/:p1/*|*=/t2/t3")
	tF("POST", "/", 404, "")
}
