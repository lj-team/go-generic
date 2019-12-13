package proxy

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProxy(t *testing.T) {

	url := "/"

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", url, nil)

	Proxy(rw, req, nil)

	resp := rw.Result()
	if resp.StatusCode != 502 {
		t.Fatal("Not work with empty opts")
	}

	o := &Opts{
		AddOriginUrl: true,
	}

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", url, nil)

	Proxy(rw, req, o)

	resp = rw.Result()
	if resp.StatusCode != 502 {
		t.Fatal("Not work with empty opts")
	}

	o.GW = []string{"http://checkip.dyn.com"}

	rw = httptest.NewRecorder()
	req = httptest.NewRequest("GET", url, nil)

	Proxy(rw, req, o)

	resp = rw.Result()

	if resp.StatusCode != 200 {
		t.Fatal("invalid status code")
	}

	data, _ := ioutil.ReadAll(resp.Body)
	if strings.Index(string(data), "Current IP Address") < 0 {
		t.Fatal("Invalid response body")
	}
}
