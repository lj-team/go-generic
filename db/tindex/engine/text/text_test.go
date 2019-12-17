package text

import (
	"testing"
)

func TestTextEngine(t *testing.T) {

	eng := New()
	if eng == nil {
		t.Fatal("New not work")
	}

	data := `
  "Привет, друг!" hello
  "Приветики" hello
  "Как твои дела?" how_are_you
  "Сегодня отличный день. Не правда ли?" is_good_day
  "Как твои дела, друг?" how_are_you
  "Дорый день!" hello
  "Привет" hello
  `

	eng.LoadString(data)

	res := eng.Search("Как твои дела", 0.5)

	if len(res) != 2 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Как твои дела?" || res[0].Value != "how_are_you" {
		t.Fatal("Invalid first value")
	}

	if res[1].Name != "Как твои дела, друг?" || res[1].Value != "how_are_you" {
		t.Fatal("Invalid second value")
	}

	res = eng.Search("привет", 0.4)

	if len(res) != 3 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Привет" || res[0].Value != "hello" {
		t.Fatal("Invalid first value")
	}

	if res[1].Name != "Приветики" || res[1].Value != "hello" {
		t.Fatal("Invalid second value")
	}

	eng.Del("привет")

	res = eng.Search("привет", 0.4)

	if len(res) != 2 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Приветики" || res[0].Value != "hello" {
		t.Fatal("Invalid second value")
	}

	eng.Del("unknown")

	eng.Add("Приветики", "hello2")

	res = eng.Search("привет", 0.4)

	if len(res) != 2 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Приветики" || res[0].Value != "hello2" {
		t.Fatal("Invalid second value")
	}

	eng.Close()

	if len(eng.Names) != 0 || len(eng.Data) != 0 || len(eng.Trgms) != 0 {
		t.Fatal("Close not worl")
	}
}
