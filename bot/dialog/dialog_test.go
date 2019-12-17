package dialog

import (
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {

	list := make([]string, 64)
	uniqs := make(map[string]bool)

	dlg := New()
	defer dlg.Close()

	for i, _ := range list {
		list[i] = dlg.Get(strconv.Itoa(i))
		uniqs[list[i]] = true
	}

	if len(uniqs) != 64 {
		t.Fatal("uniqs not work")
	}

	for i, v := range list {

		if dlg.Get(strconv.Itoa(i)) != v {
			t.Fatal("reuse keys not work")
		}

	}
}
