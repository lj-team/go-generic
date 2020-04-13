package imap

import (
	"strconv"
	"testing"

	"github.com/lj-team/go-generic/cache/nodenum"
)

func TestMulti(t *testing.T) {

	n := &multi{
		hash: nodenum.New(8),
	}

	for i := 0; i < 8; i++ {
		n.nodes = append(n.nodes, makeNode(1024/8))
	}

	tG := func(k string, wait int64, status bool) {
		res, ok := n.Get(k)

		if status != ok {
			t.Fatalf("Get invalid status for key=%s", k)
		}

		if !ok && res != nil {
			t.Fatalf("Get invalid value for key=%s", k)
		}

		if ok {

			if res == nil || res.(int64) != wait {
				t.Fatalf("Get invalid value for key=%s", k)
			}

		}
	}

	tS := func(k string, v int64) {
		n.Set(k, v)
	}

	tD := func(k string) {
		n.Delete(k)
	}

	tF := func(k string, def int64, wait int64) {
		res, ok := n.Fetch(k, func(k string) (interface{}, bool) {
			return def, true
		})

		if !ok {
			t.Fatalf("Get invalid status for key=%s", k)
		}

		if res == nil && res.(int64) != wait {
			t.Fatalf("Fetch failed for key=%s", k)
		}
	}

	tG("1", 0, false)
	tS("1", 1)
	tG("1", 1, true)

	n.Flush()
	tG("1", 0, false)

	for i := int64(0); i < 10000; i++ {
		tS(strconv.FormatInt(i, 10), i)
		tG(strconv.FormatInt(i, 10), i, true)
	}

	tD("13")
	tG("13", 0, false)

	tF("13", 11, 11)
	tF("13", 13, 11)
}
