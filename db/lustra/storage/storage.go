package storage

import (
	"errors"
	"strconv"
)

type Engine interface {
	Get(string) string
	Set(string, string)
	SetNX(string, string) string
	IncBy(string, uint64) string
	DecBy(string, uint64) string
	CBAdd(string, string, int)
	HGet(string, string) string
	HSet(string, map[string]string)
	HSetNX(string, string, string) string
	HIncBy(string, string, uint64) string
	HDecBy(string, string, uint64) string
	HSetIfLess(string, string, uint64) string
	HSetIfMore(string, string, uint64) string
	SetIfLess(string, uint64) string
	SetIfMore(string, uint64) string
	Close()
}

type OpenFunc func(string) (Engine, error)

type Storage struct {
	dbh Engine
}

var (
	//	dbh     Engine              = nil
	drivers map[string]OpenFunc = map[string]OpenFunc{}
)

func Register(driver string, fn OpenFunc) {
	drivers[driver] = fn
}

func Open(driver string, dsn string) (*Storage, error) {

	if fn, has := drivers[driver]; has {
		res, err := fn(dsn)

		if err != nil {
			return nil, err
		}

		st := &Storage{dbh: res}
		return st, nil
	}

	return nil, errors.New("unknown driver")
}

func (s *Storage) Close() {
	if s.dbh != nil {
		s.dbh.Close()
		s.dbh = nil
	}
}

func (s *Storage) Get(k string) string {
	if s.dbh == nil {
		return ""
	}
	return s.dbh.Get(k)
}

func (s *Storage) Set(k, v string) {
	if s.dbh != nil {
		s.dbh.Set(k, v)
	}
}

func (s *Storage) IncBy(k string, val uint64) string {
	if s.dbh != nil {
		return s.dbh.IncBy(k, val)
	}

	return strconv.FormatUint(val, 10)
}

func (s *Storage) DecBy(k string, val uint64) string {
	if s.dbh != nil {
		return s.dbh.DecBy(k, val)
	}

	return "0"
}

func (s *Storage) CBAdd(list string, val string, limit int) {
	if s.dbh != nil {
		s.dbh.CBAdd(list, val, limit)
	}
}

func (s *Storage) HGet(hash string, key string) string {
	if s.dbh != nil {
		return s.dbh.HGet(hash, key)
	}

	return ""
}

func (s *Storage) HSet(hash string, pairs map[string]string) {
	if s.dbh != nil {
		s.dbh.HSet(hash, pairs)
	}
}

func (s *Storage) HIncBy(hash string, key string, inc uint64) string {

	if s.dbh != nil {
		return s.dbh.HIncBy(hash, key, inc)
	}

	return strconv.FormatUint(inc, 10)
}

func (s *Storage) HDecBy(hash string, key string, dec uint64) string {
	if s.dbh != nil {
		return s.dbh.HDecBy(hash, key, dec)
	}

	return "0"
}

func (s *Storage) SetNX(k, v string) string {
	if s.dbh != nil {
		return s.dbh.SetNX(k, v)
	}

	return v
}

func (s *Storage) HSetNX(h, k, v string) string {
	if s.dbh != nil {
		return s.dbh.HSetNX(h, k, v)
	}

	return v
}

func (s *Storage) SetIfMore(k string, v uint64) string {
	if s.dbh != nil {
		return s.dbh.SetIfMore(k, v)
	}

	return strconv.FormatUint(v, 10)
}

func (s *Storage) SetIfLess(k string, v uint64) string {
	if s.dbh != nil {
		return s.dbh.SetIfLess(k, v)
	}

	return strconv.FormatUint(v, 10)
}

func (s *Storage) HSetIfMore(h string, k string, v uint64) string {
	if s.dbh != nil {
		return s.dbh.HSetIfMore(h, k, v)
	}

	return strconv.FormatUint(v, 10)
}

func (s *Storage) HSetIfLess(h string, k string, v uint64) string {
	if s.dbh != nil {
		return s.dbh.HSetIfLess(h, k, v)
	}

	return strconv.FormatUint(v, 10)
}
