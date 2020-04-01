package dialog

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/lj-team/go-generic/cache"
)

var (
	nextId uint64
	ch     *cache.Cache
)

func init() {

	nextId = uint64(time.Now().UnixNano())
	ch = cache.New("size=64000 nodes=16 ttl=86400")
}

func Get(key string) string {

	res := ch.Fetch(key, func(k string) interface{} {
		val := atomic.AddUint64(&nextId, 1)
		return fmt.Sprintf("%016x", val)
	})

	return res.(string)
}
