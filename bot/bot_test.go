package bot

import (
	"strings"
	"testing"
)

func TestBot(t *testing.T) {

	bot, _ := New("bot", strings.NewReader(""), strings.NewReader(`123:
  - 123
`), 0.1)

	bot.AddRule("123", "123")

	msg, ok := bot.Message("1", "123")

	if msg != "123" || !ok {
		t.Fatal("answer func not work")
	}

	msg, ok = bot.Message("1", "12")
	if msg != "" || ok {
		t.Fatal("failed")
	}

	msg, ok = bot.Message("1", "")
	if msg != "" || ok {
		t.Fatal("failed")
	}
}

func TestPrepapreQuestion(t *testing.T) {

	data := map[string]string{
		"Привет, мой друг!":                  "привет,мой друг!",
		"   Это супер игра. Правда что ли?!": "это супер игра.правда что ли?!",
	}

	for k, v := range data {

		if prepareQuestion(k) != v {
			t.Fatal("prepareQuestion failed " + k)
		}

	}
}
