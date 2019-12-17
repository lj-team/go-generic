package dict

import (
	"strings"
	"testing"
)

func TestDict(t *testing.T) {

	srs := `
123 12
124 12
131 13
`

	load(strings.NewReader(srs))

	if v, _ := Get("123"); v != "12" {
		t.Fatal("123 failed")
	}

	if v, _ := Get("124"); v != "12" {
		t.Fatal("124 failed")
	}

	if _, h := Get("125"); h {
		t.Fatal("125 failed")
	}
}
