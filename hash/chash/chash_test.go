package chash

import (
	"strconv"
	"testing"
)

func TestConsistent(t *testing.T) {
	hash := New(10)

	wait := []int{9, 5, 1, 4, 9, 5, 1, 4, 9, 5}

	for i := int64(0); i < 10; i++ {
		key := []byte(strconv.FormatInt(i, 10))
		if hash.Get(key) != wait[int(i)] {
			t.Fatal("invalid function value")
		}
	}

	hnul := New(1)

	if hnul.Get([]byte("123")) != 0 {
		t.Fatal("invalid return for 1l value")
	}

	if hnul.Next(0) != 0 {
		t.Fatal("invalid Next value for 1 value")
	}
}
