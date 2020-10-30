package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lj-team/go-generic/db/lustra/server/handler"
	"github.com/lj-team/go-generic/db/lustra/storage"
	"github.com/lj-team/go-generic/log"
	"github.com/lj-team/go-generic/net/http/server"
	ul "github.com/lj-team/go-generic/net/http/url/params"
	"github.com/lj-team/go-generic/net/udp"
)

func Start(addr string, driver string, opts string) error {

	st, err := storage.Open(driver, opts)
	if err != nil {
		return err
	}

	return udp.Server(addr, func(in []byte) []byte {
		return handler.Handler(st, in)
	})
}

func BatchServer(udpAddr string, httpAddr string, driver string, opts string) {

	defer func() {

		if r := recover(); r != nil {
			log.Error(fmt.Sprint(r))
			log.Stack("error")
		}

	}()

	st, err := storage.Open(driver, opts)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := udp.Server(udpAddr, func(in []byte) []byte {
			return handler.Handler(st, in)
		}); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	batchHandler := func(rw http.ResponseWriter, req *http.Request, prms ul.Params) {

		if data, err := ioutil.ReadAll(req.Body); err == nil {
			go handler.Batch(st, bytes.NewReader(data))
		}

		rw.Header().Set("Context-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(200)
		rw.Write([]byte("OK"))
	}

	server.RegHandler("POST", "/batch", batchHandler)

	if err := server.Start(httpAddr); err != nil {
		panic(err)
	}
}
