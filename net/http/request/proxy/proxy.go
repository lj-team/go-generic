package proxy

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lj-team/go-generic/log"
	"github.com/lj-team/go-generic/net/http/errors"
)

type Opts struct {
	GW           []string
	AddOriginUrl bool
	Timeout      time.Duration
}

func Proxy(rw http.ResponseWriter, req *http.Request, o *Opts) {

	if o == nil || len(o.GW) == 0 {
		errors.Send(rw, 502)
		return
	}

	if o.Timeout == 0 {
		o.Timeout = time.Second * 5
	}

	for _, gw := range o.GW {

		url := gw

		if o.AddOriginUrl {
			url += req.RequestURI
		}

		content, _ := ioutil.ReadAll(req.Body)

		var input io.Reader

		if len(content) > 0 {
			input = bytes.NewReader(content)
		}

		rn, err := http.NewRequest(req.Method, url, input)
		if err != nil {
			continue
		}

		rn.Header = req.Header

		tr := &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}

		ua := http.Client{
			Timeout:   o.Timeout,
			Transport: tr,
		}

		res, e := ua.Do(rn)
		if e != nil {
			log.Logger("error", e.Error())
			continue
		}

		content, _ = ioutil.ReadAll(res.Body)

		for h := range res.Header {
			rw.Header().Set(h, res.Header.Get(h))
		}

		rw.WriteHeader(res.StatusCode)

		if len(content) > 0 {
			rw.Write(content)
		}

		return
	}

	errors.Send(rw, 504)
}
