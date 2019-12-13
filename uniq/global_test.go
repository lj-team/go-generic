package uniq

import (
	"testing"
)

func TestGlobal(t *testing.T) {

	for i := 0; i < 10; i++ {

		val := Next()
		if val == "" {
			t.Fatal("Next not work")
		}

		if !Check(val, true) {
			t.Fatal("Check not work")
		}

	}

}
