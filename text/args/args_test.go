package args

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lj-team/go-generic/text"
)

func TestTextToken(t *testing.T) {

	tP := func(str string, wait []string) {
		res := Parse(str)

		if wait == nil && res != nil {
			t.Fatal("Error parse: " + str)
		}

		if res == nil && wait != nil {
			t.Fatal("Error parse: " + str)
		}

		if !text.CmpSlices(res, wait) {
			t.Fatal("Error parse: " + str)
		}

		res = Read(strings.NewReader(str))

		if wait == nil && res != nil {
			t.Fatal("Error parse: " + str)
		}

		if res == nil && wait != nil {
			t.Fatal("Error parse: " + str)
		}

		if !text.CmpSlices(res, wait) {
			t.Fatal("Error parse: " + str)
		}
	}

	tP(`""`, []string{""})
	tP("123", []string{"123"})
	tP(" 123", []string{"123"})
	tP("123\\s123", []string{"123 123"})
	tP("123\\ 123", []string{"123 123"})
	tP("535 t456 876l", []string{"535", "t456", "876l"})
	tP("welcome to \"4th test\"", []string{"welcome", "to", "4th test"})
	tP("welcome to \"5th test\" 098765", []string{"welcome", "to", "5th test", "098765"})
	tP("welcome to \"6th test", nil)
	tP("welco\\me to \"7th test\"", nil)
	tP("только тест", []string{"только", "тест"})
	tP("\"\\\\\"", []string{"\\"})

	tE := func(data []string, wait string) {
		res := ToString(data...)

		if res != wait {
			t.Fatal(fmt.Sprint("Error ecnode: ", data))
		}

		nd := Parse(res)
		if !text.CmpSlices(data, nd) {
			t.Fatal(fmt.Sprint("Error back convert: ", data))
		}
	}

	tE([]string{""}, `""`)
	tE([]string{"1", "2", "3"}, "1 2 3")
	tE([]string{"1 2", "3"}, `"1 2" 3`)
	tE([]string{"0", "1 2", "3"}, `0 "1 2" 3`)
	tE([]string{"0\n"}, `"0\n"`)
	tE([]string{`1\2`}, `"1\\2"`)
	tE([]string{`"Тест"`, "игра"}, `"\"Тест\"" игра`)

	res := make([]string, 0, 24)
	r := ParseTo("welcome to \"5th test\" 098765", res)
	if cap(r) != 24 {
		t.Fatal("Invalid list cap")
	}

	if strings.Join(r, "|") != strings.Join([]string{"welcome", "to", "5th test", "098765"}, "|") {
		t.Fatal("ParseTo failed")
	}

	r = ParseTo("welcome to \"5th test\" 098765", res)
	if cap(r) != 24 {
		t.Fatal("Invalid list cap")
	}

	if strings.Join(r, "|") != strings.Join([]string{"welcome", "to", "5th test", "098765"}, "|") {
		t.Fatal("ParseTo failed")
	}
}
