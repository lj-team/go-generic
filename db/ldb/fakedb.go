package ldb

import (
	"bytes"
	"container/list"
)

type FakeDB struct {
	lst *list.List
}

func NewFakeDB() *FakeDB {
	return &FakeDB{lst: list.New()}
}

func (db *FakeDB) Close() {

}

func (db *FakeDB) Set(key []byte, value []byte) {

	if db == nil {
		return
	}

	if value == nil || len(value) == 0 {
		db.Del(key)
		return
	}

	for e := db.lst.Front(); e != nil; e = e.Next() {
		v := e.Value.([][]byte)
		cmp := bytes.Compare(v[0], key)
		if cmp == 0 {
			v[1] = value
			return
		} else if cmp > 0 {
			db.lst.InsertBefore([][]byte{key, value}, e)
			return
		}
	}

	db.lst.PushBack([][]byte{key, value})
}

func (db *FakeDB) Get(key []byte) []byte {

	if db == nil {
		return nil
	}

	for e := db.lst.Front(); e != nil; e = e.Next() {
		v := e.Value.([][]byte)
		cmp := bytes.Compare(v[0], key)
		if cmp == 0 {
			return v[1]
		}
	}

	return nil
}

func (db *FakeDB) Has(key []byte) bool {

	return Get(key) != nil
}

func (db *FakeDB) Del(key []byte) {

	if db == nil {
		return
	}

	for e := db.lst.Front(); e != nil; e = e.Next() {
		v := e.Value.([][]byte)
		cmp := bytes.Compare(v[0], key)
		if cmp == 0 {
			db.lst.Remove(e)
			return
		}
	}
}

func (db *FakeDB) Total(prefix []byte) int64 {

	if db == nil {
		return 0
	}

	total := int64(0)

	db.ForEach(prefix, false, func(key []byte, value []byte) bool {
		total++
		return true
	})

	return total
}

func (db *FakeDB) ForEach(prefix []byte, RemovePrefix bool, fn FOR_EACH_FUNC) {

	if db == nil {
		return
	}

	for e := db.lst.Front(); e != nil; e = e.Next() {

		v := e.Value.([][]byte)
		key := v[0]

		if len(key) < len(prefix) {
			continue
		}

		cur := key[:len(prefix)]

		if bytes.Compare(cur, prefix) == 0 {
			exit := false
			if RemovePrefix {
				if len(key) == len(prefix) {
					exit = !fn([]byte{}, v[1])
				} else {
					exit = !fn(key[len(prefix):], v[1])
				}
			} else {
				exit = !fn(v[0], v[1])
			}
			if exit {
				return
			}
		}
	}
}

func (db *FakeDB) ForEachKey(prefix []byte, limit int, offset int, RemovePrefix bool, fn FOR_EACH_KEY_FUNC) {

	if limit < 1 {
		return
	}

	pos := 0
	total := 0

	db.ForEach(prefix, RemovePrefix, func(key []byte, value []byte) bool {
		pos++
		if pos >= offset {
			total++
			fn(key)
		}

		return total < limit
	})

}

func (db *FakeDB) List(prefix []byte, limit int, offset int, RemovePrefix bool) [][]byte {

	if limit < 1 {
		return [][]byte{}
	}

	pos := 0
	total := 0

	res := make([][]byte, 0, limit)

	db.ForEach(prefix, RemovePrefix, func(key []byte, value []byte) bool {
		pos++
		if pos >= offset {
			total++

			l1 := make([]byte, len(key))
			copy(l1, key)

			res = append(res, l1)
		}

		return total < limit
	})

	return res
}
