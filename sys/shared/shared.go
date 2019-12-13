package shared

import (
	"hash/crc32"
	"sync"
)

type Shared struct {
	size int
	list []*sync.RWMutex
}

func New(size int) *Shared {
	r := &Shared{
		size: size,
		list: make([]*sync.RWMutex, size),
	}

	for i, _ := range r.list {
		r.list[i] = &sync.RWMutex{}
	}

	return r
}

func (s *Shared) node(key []byte) int {
	return int(crc32.ChecksumIEEE(key) % uint32(s.size))
}

func (s *Shared) RLock(key []byte) {
	s.list[s.node(key)].RLock()
}

func (s *Shared) RUnLock(key []byte) {
	s.list[s.node(key)].RUnlock()
}

func (s *Shared) Lock(key []byte) {
	s.list[s.node(key)].Lock()
}

func (s *Shared) UnLock(key []byte) {
	s.list[s.node(key)].Unlock()
}
