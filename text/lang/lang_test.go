package lang

import (
	"testing"
)

func TestLang(t *testing.T) {

	rus := []string{
		"Привет, друг!",
		"Какой чудестный день.",
	}

	eng := []string{
		"Hello, world.",
	}

	for _, v := range rus {
		if !IsRus(v) {
			t.Fatal("IsRus failed " + v)
		}

		if IsEng(v) {
			t.Fatal("IsEng failed " + v)
		}
	}

	for _, v := range eng {
		if IsRus(v) {
			t.Fatal("IsRus failed " + v)
		}

		if !IsEng(v) {
			t.Fatal("IsEng failed " + v)
		}
	}

	cyr := []string{
		"привет, мир!",
		"только тест.",
		"але в той же час збільшує вірогідність",
	}

	for _, v := range cyr {
		if !IsCyr(v) {
			t.Fatal("IsCyr not work")
		}
	}

	noncyr := []string{
		"test and tets",
	}

	for _, v := range noncyr {
		if IsCyr(v) {
			t.Fatal("IsCyr not work")
		}
	}

	if !IsRusWord("игра") {
		t.Fatal("IsRuWord not work")
	}

	if IsRusWord("еуыеt") {
		t.Fatal("IsRuWord not work")
	}

	if !IsEngWord("game") {
		t.Fatal("IsEngWord not work")
	}

	if IsEngWord("конец") {
		t.Fatal("IsEngWord not work")
	}
}
