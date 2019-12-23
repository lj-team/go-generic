package ldb

import (
	"testing"

	"github.com/lj-team/go-generic/db/ldb"
	"github.com/lj-team/go-generic/sys/shared"
)

func TestLDB(t *testing.T) {

	dbh, _ := ldb.Open("test=1")

	en := &engine{
		dbh:    dbh,
		shared: shared.New(64),
	}

	defer en.Close()

	en.Set("1", "1")
	en.Set("2", "2")
	en.Set("3", "3")
	en.Set("4", "4")
	en.Set("5", "5")

	en.Set("1", "")

	tf := func(k, v string) {
		res := en.Get(k)
		if res != v {
			t.Fatal("failed for: " + string(k))
		}
	}

	tf("1", "")
	tf("2", "2")
	tf("5", "5")

	tInc := func(k string, val uint64, wait string) {
		res := en.IncBy(k, val)
		if res != wait {
			t.Fatal("IncBy failed: k=", k, " val=", val, " wait=", wait)
		}
	}

	tInc("1", 2, "2")
	tInc("6", 0, "0")
	tInc("5", 7, "12")
	tInc("6", 1, "1")
	tInc("5", 3, "15")
	tInc("8", 100000, "100000")

	tDec := func(k string, val uint64, wait string) {
		res := en.DecBy(k, val)
		if res != wait {
			t.Fatal("DecBy failed: k=", k, " val=", val, " wait=", wait)
		}
	}

	tDec("7", 1, "0")
	tDec("6", 2, "0")
	tDec("5", 2, "13")
	tDec("5", 4, "9")

	tCBA := func(k string, value string, limit int, wait string) {
		en.CBAdd(k, value, limit)
		res := en.Get(k)
		if res != wait {
			t.Fatal("expect: " + wait)
		}
	}

	tCBA("5", "123", 5, `["123"]`)
	tCBA("9", "1", 5, `["1"]`)
	tCBA("8", "44", 5, `["44"]`)
	tCBA("5", "1", 5, `["123","1"]`)
	tCBA("5", "12", 5, `["123","1","12"]`)
	tCBA("5", "123", 5, `["123","1","12","123"]`)
	tCBA("5", "1234", 5, `["123","1","12","123","1234"]`)
	tCBA("5", "12345", 5, `["1","12","123","1234","12345"]`)
	tCBA("5", "123456", 3, `["1234","12345","123456"]`)
	tCBA("5", "12345", 0, ``)

	tHSet := func(h, k, v string) {
		en.HSet(h, map[string]string{k: v})

		if en.HGet(h, k) != v {
			t.Fatal("failed for: " + h + "." + k)
		}
	}

	tHSet("h1", "1", "2")
	tHSet("h1", "2", "3")
	tHSet("h1", "3", "4")
	tHSet("h1", "4", "5")
	tHSet("h1", "5", "6")
	tHSet("h1", "6", "7")
	tHSet("h1", "2", "")

	tHInc := func(h string, k string, val uint64, wait string) {
		res := en.HIncBy(h, k, val)
		if res != wait {
			t.Fatal("HIncBy failed: h=", h, " k=", k, " v=", val, " wait=", wait)
		}
	}

	tHInc("h1", "2", 3, "3")
	tHInc("h1", "3", 1, "5")
	tHInc("h1", "0", 0, "0")
	tHInc("h1", "6", 6, "13")

	tHDec := func(h string, k string, val uint64, wait string) {
		res := en.HDecBy(h, k, val)
		if res != wait {
			t.Fatal("HDecBy failed: h=", h, " k=", k, " v=", val, " wait=", wait)
		}
	}

	tHDec("h1", "0", 1, "0")
	tHDec("h1", "6", 3, "10")
	tHDec("h1", "6", 4, "6")
	tHDec("h1", "3", 5, "0")
	tHDec("h1", "2", 7, "0")

	tHSetNX := func(h, k, v, w string) {
		res := en.HSetNX(h, k, v)
		if res != w {
			t.Fatal("HSetNX not work for h=", h, " k=", k, " v=", v, " w=", w)
		}
		if en.HGet(h, k) != w {
			t.Fatal("HSetNX not store data for h=", h, " k=", k, " v=", v, " w=", w)
		}
	}

	tHSetNX("h1", "15", "20", "20")
	tHSetNX("h1", "15", "21", "20")
	tHSetNX("h1", "0", "18", "18")
	tHSetNX("h1", "0", "", "18")

	tSetNX := func(k, v, w string) {
		res := en.SetNX(k, v)
		if res != w {
			t.Fatal("SetNX not work for k=", k, " v=", v, " w=", w)
		}
		if en.Get(k) != w {
			t.Fatal("SetNX not store data for k=", k, " v=", v, " w=", w)
		}
	}

	tSetNX("sn", "1", "1")
	tSetNX("sn", "2", "1")
	tSetNX("5", "19", "19")
	tSetNX("5", "", "19")

	tSetIfMore := func(k string, v uint64, w string) {
		res := en.SetIfMore(k, v)
		if res != w {
			t.Fatalf("SetIfMore faiiled k=%s d=%d w=%s", k, v, w)
		}
	}

	tSetIfMore("5", 8, "19")
	tSetIfMore("5", 15, "19")
	tSetIfMore("5", 20, "20")
	tSetIfMore("sif", 3, "3")
	tSetIfMore("sif1", 0, "")

	tSetIfLess := func(k string, v uint64, w string) {
		res := en.SetIfLess(k, v)
		if res != w {
			t.Fatalf("SetIfLess failed k=%s d=%d w=%s", k, v, w)
		}
	}

	tSetIfLess("sif", 5, "3")
	tSetIfLess("sif", 2, "2")
	tSetIfLess("sif", 0, "")
	tSetIfLess("sil", 13, "13")

	tHSetIfMore := func(h string, k string, v uint64, w string) {
		res := en.HSetIfMore(h, k, v)
		if res != w {
			t.Fatalf("HSetIfMore faiiled h=%s k=%s d=%d w=%s", h, k, v, w)
		}
	}

	tHSetIfMore("h1", "15", 13, "20")
	tHSetIfMore("h1", "15", 0, "20")
	tHSetIfMore("h1", "15", 22, "22")
	tHSetIfMore("h1", "19", 0, "")
	tHSetIfMore("h1", "19", 1, "1")

	tHSetIfLess := func(h string, k string, v uint64, w string) {
		res := en.HSetIfLess(h, k, v)
		if res != w {
			t.Fatalf("HSetIfLess faiiled h=%s k=%s d=%d w=%s", h, k, v, w)
		}
	}

	tHSetIfLess("h1", "19", 30, "1")
	tHSetIfLess("h1", "19", 0, "")
	tHSetIfLess("h1", "15", 12, "12")
	tHSetIfLess("h1", "30", 100, "100")
	tHSetIfLess("h1", "30", 90, "90")
}
