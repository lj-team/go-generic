package dialog

import (
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {

	data := map[string]string{}

	tN := func(k string) {

		val := Get(k)
		if Get(k) != val {
			t.Fatalf("Get caching failed for: %s", k)
		}
		data[k] = val
	}

	for i := 0; i < 100; i++ {
		tN(strconv.Itoa(i))
	}

	for k, v := range data {
		if Get(k) != v {
			t.Fatalf("Caching failed for: %s", k)
		}
	}
}
