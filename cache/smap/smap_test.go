package smap

import (
	"testing"
)

func TestDMap(t *testing.T) {

	tN := func(dsn string, ok bool) {
		n, err := New(dsn)
		if ok && err != nil || !ok && err == nil {
			t.Fatalf("New failed for: %s", dsn)
		}
		if n == nil && ok || n != nil && !ok {
			t.Fatalf("New return invalid value: %s", dsn)
		}
	}

	tN("size=128", true)
	tN("test=t=rue", false)
	tN("size=1024 nodes=4", true)
	tN("size=12", true)
}
