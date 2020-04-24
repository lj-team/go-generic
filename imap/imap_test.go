package imap

import (
	"testing"
)

func TestIMap(t *testing.T) {

	m := New()
	if m == nil {
		t.Fatal("New failed")
	}

	tS := func(k int64, v string, has bool, size int) {
		m.Set(k, v)

		res, ok := m.Get(k)
		if ok != has {
			t.Fatalf("Get flag invalid valid for: k=%d v=%s", k, v)
		}

		if res != v {
			t.Fatalf("Get return invalid value for: k=%d v=%s", k, v)
		}

		if m.Has(k) != has {
			t.Fatalf("Has failed for: k=%d v=%s", k, v)
		}

		if m.Len() != size {
			t.Fatalf("Len failed for: k=%d v=%s", k, v)
		}
	}

	tS(1, "12", true, 1)
	tS(2, "23", true, 2)
	tS(3, "34", true, 3)
	tS(4, "45", true, 4)
	tS(2, "", false, 3)

	m.Delete(4)
	if m.Has(4) || m.Len() != 2 {
		t.Fatal("Delete failed")
	}

	cnt := 0

	m.ForEach(func(k int64, v string) {
		cnt++
		if v == "" {
			t.Fatal("ForEach empty string")
		}
		if k != 1 && k != 3 {
			t.Fatal("ForEach unknown key")
		}

		if v != "12" && v != "34" {
			t.Fatal("ForEach unknown value")
		}
	})

	if cnt != 2 {
		t.Fatal("ForEach invalid iteration counter")
	}

	m.ForEach(nil)

	m.Reset()
	if m.Len() != 0 || m.Has(1) {
		t.Fatal("Reset failed")
	}
}
