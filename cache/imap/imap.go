package imap

import (
	"errors"

	"github.com/lj-team/go-generic/cache/nodenum"
	kvs "github.com/lj-team/go-generic/text/kv"
)

type FetchFunc func(string) (interface{}, bool)

type IMap interface {
	Get(key string) (interface{}, bool)

	Set(key string, val interface{})

	Fetch(key string, fn FetchFunc) (interface{}, bool)

	Delete(key string)

	Flush()
}

func New(dsn string) (IMap, error) {

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
