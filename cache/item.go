package cache

import (
	"time"
)

type item struct {
	data   interface{} // value
	expire int64       // expire time
}

func (i *item) isAlive() bool {
	return time.Now().Unix() < i.expire
}
