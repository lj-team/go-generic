package imap

import (
	"github.com/lj-team/go-generic/cache/nodenum"
)

type multi struct {
	hash  *nodenum.NodeNum
	nodes []*node
}

func (m *multi) Get(key string) (interface{}, bool) {
	num := m.hash.Get(key)
	return m.nodes[num].Get(key)
}

func (m *multi) Set(key string, val interface{}) {
	num := m.hash.Get(key)
	m.nodes[num].Set(key, val)
}

func (m *multi) Fetch(key string, fn FetchFunc) (interface{}, bool) {
	num := m.hash.Get(key)
	return m.nodes[num].Fetch(key, fn)
}

func (m *multi) Delete(key string) {
	num := m.hash.Get(key)
	m.nodes[num].Delete(key)
}

func (m *multi) Flush() {
	for _, n := range m.nodes {
		n.Flush()
	}
}
