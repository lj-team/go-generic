package ldb

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/lj-team/go-generic/log"
)

var store Storage = nil
var proxyConfig *Config

func Init(cfg *Config) {
	if store == nil {
		New(cfg, true)
	}
}

func TestInit() {
	Open("test=true")
}

func Close() {
	if store != nil {
		store.Close()
	}
	store = nil
}

func InitProxy(cfg *Config) {
	proxyConfig = cfg
}

func Use(base string) {

	cfg := &Config{
		Path:        proxyConfig.Path + "/" + base,
		FileSize:    proxyConfig.FileSize,
		Compression: proxyConfig.Compression,
		ReadOnly:    proxyConfig.ReadOnly,
	}

	Close()
	Init(cfg)
}

func Set(key []byte, value []byte) {
	store.Set(key, value)
}

func Get(key []byte) []byte {
	return store.Get(key)
}

func Has(key []byte) bool {
	return store.Has(key)
}

func Del(key []byte) {
	store.Del(key)
}

func Total(prefix []byte) int64 {
	return store.Total(prefix)
}

func List(prefix []byte, limit int, offset int, RemovePrefix bool) [][]byte {

	return store.List(prefix, limit, offset, RemovePrefix)
}

func ForEach(prefix []byte, RemovePrefix bool, fn FOR_EACH_FUNC) {

	store.ForEach(prefix, RemovePrefix, fn)
}

func ForEachKey(prefix []byte, limit int, offset int, RemovePrefix bool, fn FOR_EACH_KEY_FUNC) {

	store.ForEachKey(prefix, limit, offset, RemovePrefix, fn)
}

func TmpCopy(src string, fn func(Storage)) {

	tmpCopy := func() bool {

		defer time.Sleep(time.Second * 3)

		ts := time.Now().Unix()

		tmpName := fmt.Sprintf("%s-%d-%d", src, ts, time.Now().UnixNano()%1000)

		cmd := fmt.Sprintf("cp -r %s %s", src, tmpName)

		exec.Command("sh", "-c", cmd).Run()

		defer func() {
			exec.Command("sh", "-c", "rm -rf "+tmpName).Run()
		}()

		db, err := Open(fmt.Sprintf("path=" + tmpName))
		if err != nil {
			log.Error(err.Error() + " " + tmpName)
			return false
		}
		defer db.Close()

		fn(db)

		return true
	}

	for i := 0; i < 5; i++ {
		if tmpCopy() {
			return
		}
	}
}
