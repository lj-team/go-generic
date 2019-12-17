package dialog

import (
	"fmt"
	"time"

	ch "github.com/lj-team/go-generic/cache"
)

type Dialog struct {
	idChan chan string
	cache  *ch.Cache
	finish chan bool
}

func New() *Dialog {

	dlg := &Dialog{
		idChan: make(chan string, 128),
		cache:  ch.New("size=64000 nodes=16 ttl=86400"),
		finish: make(chan bool),
	}

	go func() {

		num := time.Now().Unix()

		for {

			num = (num + 1) & 0xffff

			var val int64 = (time.Now().Unix() & 0xffffffff) | (num << 32)

			select {

			case <-dlg.finish:

				close(dlg.idChan)
				return

			case dlg.idChan <- fmt.Sprintf("%012X", val):

			}
		}

	}()

	return dlg
}

func (dlg *Dialog) Get(key string) string {

	res := dlg.cache.Fetch(key, func(k string) interface{} {
		v := <-dlg.idChan
		return v
	}).(string)

	dlg.cache.Set(key, res)

	return res
}

func (dlg *Dialog) Close() {
	close(dlg.finish)
	for _ = range dlg.idChan {
	}
}
