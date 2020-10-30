package lustra

import (
	"github.com/lj-team/go-generic/db/lustra/server"

	_ "github.com/lj-team/go-generic/db/lustra/storage/engine/cache"
	_ "github.com/lj-team/go-generic/db/lustra/storage/engine/ldb"
)

func Start(addr string, driver string, opts string) error {

	if err := server.Start(addr, driver, opts); err != nil {
		return err
	}

	return nil
}

func BatchServer(udpAddr string, httpAddr string, driver string, opts string) {
	server.BatchServer(udpAddr, httpAddr, driver, opts)
}
