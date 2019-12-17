package bot

import (
	"testing"
)

func TestBot(t *testing.T) {

	msg := ""
	did := ""

	bot := New(&Config{RulesFile: "", AnswerFunc: func(id, in string) (string, bool) {

		if did != "" && did != id {
			t.Fatal("invalid did")
		}

		did = id
		msg = in
		return "", false
	}, History: ""})

	bot.AddRule("123", "123")

	bot.Message("1", "123")

	if msg != "123" {
		t.Fatal("answer func not work")
	}

	if did == "" {
		t.Fatal("did not set")
	}

	bot.Message("1", "12")
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
