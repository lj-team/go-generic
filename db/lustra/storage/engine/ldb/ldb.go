package ldb

import (
	"encoding/json"
	"strconv"

	"github.com/lj-team/go-generic/db/ldb"
	"github.com/lj-team/go-generic/db/lustra/storage"
	"github.com/lj-team/go-generic/sys/shared"
)

type engine struct {
	dbh    ldb.Storage
	shared *shared.Shared
}

func (s *engine) Close() {
	s.dbh.Close()
}

func (s *engine) Get(k string) string {
	return string(s.dbh.Get([]byte(k)))
}

func (s *engine) Set(k, v string) {

	key := []byte(k)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	if len(v) == 0 {
		s.dbh.Del(key)
	} else {
		s.dbh.Set(key, []byte(v))
	}

}

func (s *engine) SetNX(k, v string) string {

	key := []byte(k)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	val := s.Get(k)
	if val != "" || v == "" {
		return val
	}

	s.dbh.Set(key, []byte(v))

	return v
}

func (s *engine) SetIfMore(k string, v uint64) string {
	key := []byte(k)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	val := s.Get(k)
	old, _ := strconv.ParseUint(val, 10, 64)

	if old < v {
		val = strconv.FormatUint(v, 10)
		s.dbh.Set(key, []byte(val))
		return val
	} else if old == 0 {
		s.dbh.Del(key)
		val = ""
	}

	return val
}

func (s *engine) SetIfLess(k string, v uint64) string {
	key := []byte(k)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	if v == 0 {
		s.dbh.Del(key)
		return ""
	}

	val := s.Get(k)
	if val == "" {
		val = strconv.FormatUint(v, 10)
		s.dbh.Set(key, []byte(val))
		return val
	}

	old, _ := strconv.ParseUint(val, 10, 64)
	if old > v {
		val = strconv.FormatUint(v, 10)
		s.dbh.Set(key, []byte(val))
		return val
	}

	return val
}

func (s *engine) IncBy(k string, cnt uint64) string {

	key := []byte(k)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	strVal := string(s.dbh.Get(key))

	val, _ := strconv.ParseUint(strVal, 10, 64)
	val += cnt

	if val == 0 {
		s.dbh.Del(key)
		return "0"
	}
	strVal = strconv.FormatUint(val, 10)
	s.dbh.Set(key, []byte(strVal))
	return strVal
}

func (s *engine) DecBy(k string, cnt uint64) string {

	key := []byte(k)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	strVal := string(s.dbh.Get(key))

	val, _ := strconv.ParseUint(strVal, 10, 64)

	if val <= cnt {
		val = 0
	} else {
		val -= cnt
	}

	if val == 0 {
		s.dbh.Del(key)
		return "0"
	}
	strVal = strconv.FormatUint(val, 10)
	s.dbh.Set(key, []byte(strVal))
	return strVal
}

func (s *engine) UHeap(list string, value string, limit int) {
	key := []byte(list)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	if limit < 1 {
		s.dbh.Del(key)
		return
	}

	data := s.dbh.Get(key)

	var queue []string

	if len(data) > 3 {
		if err := json.Unmarshal(data, &queue); err == nil {

			idx := -1

			for i, v := range queue {
				if v == value {
					idx = i
					break
				}
			}

			if idx >= 0 {

				for i := idx; i < len(queue)-1; i++ {
					queue[i] = queue[i+1]
				}
				queue[len(queue)-1] = value

			} else {
				queue = append(queue, value)
			}

			if len(queue) > limit {
				queue = queue[len(queue)-limit:]
			}
			data, _ = json.Marshal(queue)
			s.dbh.Set(key, data)
			return
		}
	}

	queue = []string{value}
	data, _ = json.Marshal(queue)
	s.dbh.Set(key, data)
}

func (s *engine) CBAdd(list string, value string, limit int) {
	key := []byte(list)

	s.shared.Lock(key)
	defer s.shared.UnLock(key)

	if limit < 1 {
		s.dbh.Del(key)
		return
	}

	data := s.dbh.Get(key)

	var queue []string

	if len(data) > 3 {
		if err := json.Unmarshal(data, &queue); err == nil {
			queue = append(queue, value)
			if len(queue) > limit {
				queue = queue[len(queue)-limit:]
			}
			data, _ = json.Marshal(queue)
			s.dbh.Set(key, data)
			return
		}
	}

	queue = []string{value}
	data, _ = json.Marshal(queue)
	s.dbh.Set(key, data)
}

func (s *engine) getMap(hash []byte) map[string]string {

	data := s.dbh.Get([]byte(hash))
	var mp map[string]string

	if len(data) < 2 {
		return map[string]string{}
	}

	if err := json.Unmarshal(data, &mp); err != nil {
		return map[string]string{}
	}

	return mp
}

func (s *engine) saveMap(hash []byte, mp map[string]string) {

	if mp == nil || len(mp) == 0 {
		s.dbh.Del(hash)
	} else {
		data, _ := json.Marshal(mp)
		s.dbh.Set(hash, data)
	}
}

func (s *engine) HSet(hash string, pairs map[string]string) {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hk)
	defer s.saveMap(hk, mp)

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

	mp := s.getMap(hk)
	defer s.saveMap(hk, mp)

	has := mp[key]
	if has != "" || value == "" {
		return has
	}

	mp[key] = value

	return value
}

func (s *engine) HSetIfMore(hash string, key string, value uint64) string {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hk)
	defer s.saveMap(hk, mp)

	res := mp[key]

	old, _ := strconv.ParseUint(res, 10, 64)

	if old < value {
		res = strconv.FormatUint(value, 10)
		mp[key] = res
	} else if old == 0 {
		delete(mp, key)
		res = ""
	}

	return res
}

func (s *engine) HSetIfLess(hash string, key string, value uint64) string {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hk)
	defer s.saveMap(hk, mp)

	res, has := mp[key]

	if value == 0 {
		delete(mp, key)
		return ""
	}

	if !has {
		res = strconv.FormatUint(value, 10)
		mp[key] = res
		return res
	}

	old, _ := strconv.ParseUint(res, 10, 64)

	if old > value {
		res = strconv.FormatUint(value, 10)
		mp[key] = res
	} else if old == 0 {
		delete(mp, key)
		res = ""
	}

	return res
}

func (s *engine) HGet(hash string, key string) string {
	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hk)
	return mp[key]
}

func (s *engine) HIncBy(hash string, key string, cnt uint64) string {

	hk := []byte(hash)

	s.shared.Lock(hk)
	defer s.shared.UnLock(hk)

	mp := s.getMap(hk)
	defer s.saveMap(hk, mp)

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

	mp := s.getMap(hk)
	defer s.saveMap(hk, mp)

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

	storage.Register("ldb", func(dsn string) (storage.Engine, error) {

		dbh, err := ldb.Open(dsn)
		if err != nil {
			return nil, err
		}

		return &engine{dbh: dbh, shared: shared.New(64)}, nil

	})

}
