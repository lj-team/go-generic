package lang

import (
	"testing"
)

func TestLang(t *testing.T) {

	data := map[string]int64{
		"123":       LANG_DI,
		"привет":    LANG_RU,
		"world":     LANG_EN,
		"ვზთიკლმნო": LANG_UN,
		"прівет":    LANG_RU | LANG_UN,
	}

	for k, v := range data {
		if Lang(k) != v {
			t.Fatal("Failed " + k)
		}
	}
}
