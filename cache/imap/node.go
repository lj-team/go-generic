package imap

import (
	"sync"
)

const (
	coeff_limit     float64 = 1.2
	coeff_threshold float64 = 0.8
)

type node struct {
	sync.RWMutex
	data      map[string]interface{}
	prev      map[string]interface{}
	limit     int
	total     int
	prevTotal int
	threshold int
}

func makeNode(limit int) *node {
	return &node{
		data:      make(map[string]interface{}, int(float64(limit)*coeff_limit)),
		prev:      map[string]interface{}{},
		limit:     limit,
		threshold: int(float64(limit) * coeff_threshold),
		total:     0,
	}
}

func (n *node) Get(key string) (interface{}, bool) {
	n.RLock()
	defer n.RUnlock()

	if v, h := n.data[key]; h {
		return v, true
	}

	if v, h := n.prev[key]; h {
		return v, true
	}

	return nil, false
}

func (n *node) Set(key string, value interface{}) {

	if value == nil {
		n.Delete(key)
		return
	}

	n.Lock()
	defer n.Unlock()

	if _, has := n.data[key]; has {
		n.data[key] = value
	} else {
		n.gc()
		n.total++
		n.data[key] = value
	}
}

func (n *node) gc() {

	if n.total >= n.threshold {
		n.prev = n.data
		n.data = make(map[string]interface{}, int(float64(n.limit)*coeff_limit))
		n.prevTotal = n.total
		n.total = 0
	}
}

func (n *node) Delete(key string) {
	n.Lock()
	defer n.Unlock()

	if _, has := n.data[key]; has {
		delete(n.data, key)
		n.total--
	}

	if _, has := n.prev[key]; has {
		delete(n.prev, key)
		n.prevTotal--
	}
}

func (n *node) Size() int {
	n.RLock()
	defer n.RUnlock()

	return n.total + n.prevTotal
}

func (n *node) Flush() {
	n.Lock()
	defer n.Unlock()

	n.data = make(map[string]interface{}, int(float64(n.limit)*coeff_limit))
	n.prev = map[string]interface{}{}
	n.total = 0
	n.prevTotal = 0
}

func (n *node) Fetch(key string, fn FetchFunc) (interface{}, bool) {

	if res, has := n.Get(key); has {
		return res, has
	}

	if res, has := fn(key); has {
		n.Set(key, res)
		return res, true
	}

	return nil, false
}
