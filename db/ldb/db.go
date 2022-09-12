package ldb

import (
	"errors"

	"github.com/lj-team/go-generic/log"
	kvs "github.com/lj-team/go-generic/text/kv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DB struct {
	ldb *leveldb.DB
}

func Open(dsn string) (Storage, error) {

	kv, err := kvs.New(dsn)
	if err != nil {
		return nil, err
	}

	if kv.GetBool("test", false) {
		if store != nil {
			store.Close()
		}
		store = NewFakeDB()
		return store, nil
	}

	path := kv.GetString("path", "")
	if path == "" {
		return nil, errors.New("path not entered")
	}

	cfg := &Config{
		Path:        kv.GetString("path", ""),
		Compression: kv.GetBool("compression", true),
		FileSize:    kv.GetInt("filesize", 16),
		ReadOnly:    kv.GetBool("readonly", false),
	}

	return New(cfg, kv.GetBool("default", false))
}

func New(cfg *Config, def bool) (*DB, error) {
	db := &DB{}

	comp := opt.NoCompression
	mul := 1
	if cfg.Compression {
		comp = opt.SnappyCompression
		mul = 2
	}

	size := cfg.FileSize * 1024 * 1024

	var err error

	log.InfoParams("Открытие базы", cfg.Path)
	db.ldb, err = leveldb.OpenFile(cfg.Path, &opt.Options{
		CompactionTableSize: size,
		WriteBuffer:         size * mul,
		Compression:         comp,
		ReadOnly:            cfg.ReadOnly,
	})

	if err != nil {
		log.ErrorParams("Ошибка открытия базы", err.Error())
		if def {
			store = nil
		}
	} else if def {
		store = db
	}

	return db, err
}

func (db *DB) Close() {
	if db != nil && db.ldb != nil {
		db.ldb.Close()
	}
}

func (db *DB) Set(key []byte, value []byte) {

	if db == nil {
		return
	}

	if value == nil || len(value) == 0 {
		db.ldb.Delete(key, nil)
	} else {
		db.ldb.Put(key, value, nil)
	}
}

func (db *DB) Get(key []byte) []byte {
	val, err := db.ldb.Get(key, nil)

	if err != nil {
		return nil
	}

	return val
}

func (db *DB) Has(key []byte) bool {
	if db == nil {
		return false
	}

	val, err := db.ldb.Has(key, nil)
	if err != nil {
		return false
	}
	return val
}

func (db *DB) Del(key []byte) {
	if db != nil {
		db.ldb.Delete(key, nil)
	}
}

func (db *DB) Total(prefix []byte) int64 {

	if db == nil {
		return 0
	}

	iter := db.ldb.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	i := int64(0)

	for iter.Next() {
		i++
	}

	return i
}

func (db *DB) List(prefix []byte, limit int, offset int, RemovePrefix bool) [][]byte {

	res := make([][]byte, 0)

	if db == nil {
		return res
	}

	iter := db.ldb.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	i := -1

	for iter.Next() {
		i++

		if i >= offset+limit {
			break
		}

		if i < offset {
			continue
		}

		var list []byte

		if RemovePrefix {
			size := len(iter.Key()) - len(prefix)
			list = make([]byte, size)
			copy(list, (iter.Key())[len(prefix):])
		} else {
			list = make([]byte, len(iter.Key()))
			copy(list, iter.Key())
		}

		res = append(res, list)
	}

	return res
}

func (db *DB) ForEach(prefix []byte, RemovePrefix bool, fn FOR_EACH_FUNC) {

	if db == nil {
		return
	}

	iter := db.ldb.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	list := make([]byte, 4096)
	var size int

	for iter.Next() {

		if RemovePrefix {
			size = len(iter.Key()) - len(prefix)
			copy(list, (iter.Key())[len(prefix):])
		} else {
			size = len(iter.Key())
			copy(list, iter.Key())
		}

		if !fn(list[:size], iter.Value()) {
			return
		}
	}
}

func (db *DB) ForEachKey(prefix []byte, limit int, offset int, RemovePrefix bool, fn FOR_EACH_KEY_FUNC) {

	if db == nil {
		return
	}

	iter := db.ldb.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	i := -1

	for iter.Next() {
		i++

		if i >= offset+limit {
			break
		}

		if i < offset {
			continue
		}

		var list []byte

		if RemovePrefix {
			size := len(iter.Key()) - len(prefix)
			list = make([]byte, size)
			copy(list, (iter.Key())[len(prefix):])
		} else {
			list = make([]byte, len(iter.Key()))
			copy(list, iter.Key())
		}

		if !fn(list) {
			return
		}
	}
}
