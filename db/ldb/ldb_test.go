package ldb

import (
	"testing"
)

func TestDB(t *testing.T) {

	Close()
	Close()

	TestInit()

	key := []byte("test.key")
	val := []byte{1, 2, 3}

	Del(key)

	res := Get(key)

	if res != nil {
		t.Fatal("expecting nil value")
	}

	if Has(key) {
		t.Fatal("Has not working")
	}

	Set(key, val)

	res = Get(key)

	if res == nil || len(res) != len(val) {
		t.Fatal("db.Get error")
	}

	if !Has(key) {
		t.Fatal("Has not work")
	}

	for i, c := range res {
		if c != val[i] {
			t.Fatal("db.Get error")
		}
	}

	res = Get(key)

	if res == nil || len(res) != len(val) {
		t.Fatal("db.Get error")
	}

	for i, c := range res {
		if c != val[i] {
			t.Fatal("db.Get error")
		}
	}

	if Total([]byte{}) != 1 {
		t.Fatal("wrong key number")
	}

	Set([]byte("tkey"), []byte("1"))

	if Total([]byte{}) != 2 {
		t.Fatal("wrong key number")
	}

	if Total([]byte("tk")) != 1 {
		t.Fatal("total prefix find error")
	}
}
