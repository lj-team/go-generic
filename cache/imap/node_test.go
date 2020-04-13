package imap

import (
	"strconv"
	"testing"
)

func TestNode(t *testing.T) {

	n := makeNode(30)
	if n == nil {
		t.Fatal("makeNode failed")
	}

	tS := func(k string, v int64) {
		n.Set(k, v)
	}

	tSize := func(wait int) {
		if n.Size() != wait {
			t.Fatalf("Size failed expect %d", wait)
		}
	}

	tG := func(k string, wait int64, found bool) {
		res, ok := n.Get(k)
		if ok != found {
			t.Fatalf("Get return invalid found flag for key=%s", k)
		}

		if ok && res == nil {
			t.Fatalf("Get return invalid found value for key=%s", k)
		}

		if ok {

			if res.(int64) != wait {
				t.Fatalf("Get return invalid found value for key=%s", k)
			}

		}
	}

	tF := func(k string, def int64, wait int64) {
		res, ok := n.Fetch(k, func(k string) (interface{}, bool) {
			return def, true
		})

		if !ok {
			t.Fatal("invalid ok")
		}

		if res.(int64) != wait {
			t.Fatalf("Fetch failed for k=%s w=%d", k, wait)
		}
	}

	tS("1", 1)
	tS("2", 2)
	tS("3", 3)

	tSize(3)

	tG("0", 0, false)
	tG("2", 2, true)

	n.Set("2", nil)

	tG("2", 0, false)

	tS("3", 33)
	tG("3", 33, true)
	tF("3", 2, 33)
	tSize(2)
	tF("2", 2, 2)

	n.Flush()
	tSize(0)
	tG("3", 0, false)

	for i := int64(0); i < 45; i++ {
		tS(strconv.FormatInt(i, 10), i)
	}

	tG("5", 5, true)

	n.Set("4", nil)

}
