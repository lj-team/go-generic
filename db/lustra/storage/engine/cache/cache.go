package cache

import (
	"encoding/json"
	"strconv"

	"github.com/lj-team/go-generic/cache/smap"
	"github.com/lj-team/go-generic/db/lustra/storage"
	"github.com/lj-team/go-generic/sys/shared"
)

type engine struct {
	cache  smap.SMap
	shared *shared.Shared
}

func (e *engine) Close() {
	e.cache.Flush()
}

func (e *engine) Set(k, v string) {

	key := []byte(k)

	e.shared.Lock(key)
	defer e.shared.UnLock(key)

	if len(v) == 0 {
		e.cache.Delete(k)
	} else {
		e.cache.Set(k, v)
	}
}

func (e *engine) SetNX(k, v string) string {

	key := []byte(k)

	e.shared.Lock(key)
	defer e.shared.UnLock(key)

	val := e.Get(k)
	if val != "" || v == "" {
		return val
	}

	e.cache.Set(k, v)
	return v
}

func (e *engine) Get(k string) string {

	res, _ := e.cache.Get(k)
	return res
}

func (e *engine) IncBy(k string, val uint64) string {

	key := []byte(k)

	e.shared.Lock(key)
	defer e.shared.UnLock(key)

	res, _ := e.cache.Get(k)
	cur, _ := strconv.ParseUint(res, 10, 64)
	cur += val

	if cur == 0 {
		e.cache.Delete(k)
		return "0"
	}

	nv := strconv.FormatUint(cur, 10)
	e.cache.Set(k, nv)
	return nv
}

func (e *engine) DecBy(k string, val uint64) string {

	key := []byte(k)

	e.shared.Lock(key)
	defer e.shared.UnLock(key)

	res, _ := e.cache.Get(k)
	cur, _ := strconv.ParseUint(res, 10, 64)

	if cur < val {
		cur = 0
	} else {
		cur -= val
	}

	if cur == 0 {
		e.cache.Delete(k)
		return "0"
	}

	nv := strconv.FormatUint(cur, 10)
	e.cache.Set(k, nv)
	return nv
}

func (s *engine) CBAdd(list string, value string, limit int) {
	key := []byte(list)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	if limit < 1 {
		s.cache.Delete(list)
		return
	}

	res, _ := s.cache.Get(list)

	var queue []string

	if len(res) > 3 {
		if err := json.Unmarshal([]byte(res), &queue); err == nil {
			queue = append(queue, value)
			if len(queue) > limit {
				queue = queue[len(queue)-limit:]
			}
			val, _ := json.Marshal(queue)
			s.cache.Set(list, string(val))
			return
		}
	}

	queue = []string{value}
	data, _ := json.Marshal(queue)
	s.cache.Set(list, string(data))
}

func (s *engine) getMap(hash string) map[string]string {

	dt, _ := s.cache.Get(hash)
	var mp map[string]string

	if len(dt) < 2 {
		return map[string]string{}
	}

	if err := json.Unmarshal([]byte(dt), &mp); err != nil {
		return map[string]string{}
	}

	return mp
}

func (s *engine) saveMap(hash string, mp map[string]string) {

	if mp == nil || len(mp) == 0 {
		s.cache.Delete(hash)
	} else {
		data, _ := json.Marshal(mp)
		s.cache.Set(hash, string(data))
	}
}

func (s *engine) HSet(hash string, pairs map[string]string) {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hash)
	defer s.saveMap(hash, mp)

	for key, value := range pairs {
		if value == "" {
			delete(mp, key)
		} else {
			mp[key] = value
		}
	}
}

func (s *engine) HSetNX(hash string, key string, value string) string {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hash)
	defer s.saveMap(hash, mp)

	has := mp[key]
	if has != "" || value == "" {
		return has
	}

	mp[key] = value
	return value
}

func (s *engine) HGet(hash string, key string) string {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hash)
	return mp[key]
}

func (s *engine) HIncBy(hash string, key string, cnt uint64) string {

	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hash)
	defer s.saveMap(hash, mp)

	valStr := mp[key]
	val, _ := strconv.ParseUint(valStr, 10, 64)

	val += cnt
	if val == 0 {
		delete(mp, key)
		return "0"
	}

	valStr = strconv.FormatUint(val, 10)
	mp[key] = valStr

	return valStr
}

func (s *engine) HDecBy(hash string, key string, cnt uint64) string {

	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hash)
	defer s.saveMap(hash, mp)

	valStr := mp[key]
	val, _ := strconv.ParseUint(valStr, 10, 64)

	if val <= cnt {
		val = 0
	} else {
		val -= cnt
	}

	if val == 0 {
		delete(mp, key)
		return "0"
	}

	valStr = strconv.FormatUint(val, 10)
	mp[key] = valStr

	return valStr
}

func init() {

	storage.Register("cache", func(dsn string) (storage.Engine, error) {

		cache, err := smap.New(dsn)
		if err == nil {
			return &engine{cache: cache, shared: shared.New(128)}, nil
		}

		return nil, err
	})

}
