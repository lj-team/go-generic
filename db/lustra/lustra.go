package lustra

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.002
// @date    2019-12-10

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
