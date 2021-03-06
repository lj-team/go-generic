package stemmer

import (
	"testing"
)

func TestTextToCode(t *testing.T) {

	data := map[string]string{
		"":                    "",
		"игра":                "игр",
		"купить ботинок":      "купить ботинок",
		"убить время":         "убить врем",
		"майорка":             "майорк",
		"пилота-разведчика":   "пилот-разведчик",
		"science":             "scienc",
		"123":                 "123",
		".":                   ".",
		"test-test":           "test-test",
		"tests-123":           "test-123",
		"москва. день 1":      "москва . день 1",
		"аквальва":            "аквалев",
		"батальи":             "баталья",
		"привет, мир!":        "привет , мир !",
		"Mission: impossible": "mission : imposs",
	}

	for k, v := range data {
		if TextToCode(k) != v {
			t.Fatal("TextToCode wait " + k + " | " + TextToCode(k))
		}
	}

}
