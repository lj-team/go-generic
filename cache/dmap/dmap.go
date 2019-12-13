package dmap

import (
	"errors"

	"github.com/lj-team/go-generic/cache/nodenum"
	kvs "github.com/lj-team/go-generic/text/kv"
)

type FetchFunc func(string) (int64, bool)

type DMap interface {
	Get(key string) (int64, bool)

	Set(key string, val int64)

	Fetch(key string, fn FetchFunc) (int64, bool)

	Delete(key string)

	IncBy(key string, inc int64) int64

	Flush()
}

func New(dsn string) (DMap, error) {

	kv, err := kvs.New(dsn)
	if err != nil {
		return nil, errors.New("invalid params")
	}

	size := kv.GetInt("size", 1024)
	if size < 32 {
		size = 32
	}

	num := kv.GetInt("nodes", 0)
	if num < 2 {
		return makeNode(size), nil
	}

	m := &multi{
		hash: nodenum.New(num),
	}

	for i := 0; i < num; i++ {
		m.nodes = append(m.nodes, makeNode(size/num))
	}

	return m, nil
}
