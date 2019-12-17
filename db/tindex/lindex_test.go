package tindex

import (
	"testing"

	_ "github.com/lj-team/go-generic/db/tindex/engine/text"
)

func TestLIndex(t *testing.T) {

	data := map[string]string{
		"Привет, друг!":  "hello",
		"Приветики":      "hello",
		"Как твои дела?": "how_are_you",
		"Сегодня отличный день. Не правда ли?": "is_good_day",
		"Как твои дела, друг?":                 "how_are_you",
		"Дорый день!":                          "hello",
		"Привет":                               "hello",
		"Полный отстой или что еще ты скажешь": "bad",
	}

	dbh, err := Open("text", "")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range data {
		dbh.Add(k, v)
	}

	res := dbh.Search("привет", 0.4)

	if len(res) != 3 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Привет" || res[0].Value != "hello" {
		t.Fatal("Invalid first value")
	}

	if res[1].Name != "Приветики" || res[1].Value != "hello" {
		t.Fatal("Invalid second value")
	}

	dbh.Del("привет")

	res = dbh.Search("привет", 0.4)

	if len(res) != 2 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Приветики" || res[0].Value != "hello" {
		t.Fatal("Invalid second value")
	}

	dbh.Del("unknown")

	dbh.Add("Приветики", "hello2")

	res = dbh.Search("привет", 0.4)

	if len(res) != 2 {
		t.Fatal("Invalid result size")
	}

	if res[0].Name != "Приветики" || res[0].Value != "hello2" {
		t.Fatal("Invalid second value")
	}
}
