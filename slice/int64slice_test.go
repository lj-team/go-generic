package slice

import (
	"testing"
)

func TestInt64Slice(t *testing.T) {

	list := []int64{5, 3, 4, 1, 2}

	sl := Int64Slice(list)
	sl.Sort()

	for i, v := range list {
		if int64(i+1) != v {
			t.Fatal("sort failed")
		}
	}
}
