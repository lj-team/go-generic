package imap

import (
	"sync"
)

type IMap struct {
	mt   sync.RWMutex
	data map[int64]string
}

func New() *IMap {
	return &IMap{
		data: make(map[int64]string),
	}
}

func (m *IMap) Get(key int64) (string, bool) {
	m.mt.RLock()
	defer m.mt.RUnlock()
	res, ok := m.data[key]
	return res, ok
}

func (m *IMap) Has(key int64) bool {
	m.mt.RLock()
	defer m.mt.RUnlock()
	_, ok := m.data[key]
	return ok
}

func (m *IMap) Set(key int64, value string) {
	m.mt.Lock()
	defer m.mt.Unlock()
	if value == "" {
		delete(m.data, key)
	} else {
		m.data[key] = value
	}
}

func (m *IMap) Delete(key int64) {
	m.Set(key, "")
}

func (m *IMap) Reset() {
	m.mt.Lock()
	defer m.mt.Unlock()
	m.data = make(map[int64]string)
}

func (m *IMap) Len() int {
	m.mt.RLock()
	defer m.mt.RUnlock()
	return len(m.data)
}

func (m *IMap) ForEach(fn func(k int64, val string)) {
	if fn == nil {
		return
	}
	m.mt.RLock()
	defer m.mt.RUnlock()
	for k, v := range m.data {
		fn(k, v)
	}
}
