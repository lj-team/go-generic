package server

import (
	"github.com/lj-team/go-generic/db/lustra/server/handler"
	"github.com/lj-team/go-generic/db/lustra/storage"
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
