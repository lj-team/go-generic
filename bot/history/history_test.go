package history

import (
	"testing"
)

func TestTextPreproc(t *testing.T) {

	data := map[string]string{
		"test":       "test",
		" test ":     "test",
		"abc\tdef\n": "abc def",
		"Это ёж! ":   "Это еж!",
		"Ёжик\r\n":   "Ежик",
	}

	for k, v := range data {

		if textPreproc(k) != v {
			t.Fatal("textPreproc " + k)
		}

	}
}
