package trgm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lj-team/go-generic/slice"
)

func TestTrgm(t *testing.T) {

	mtT := func(src string, wait string) {

		res := makeText(strings.NewReader(src))
		if res != wait {
			t.Fatal("invalid result for: " + src)
		}
	}

	mtT("", "___")
	mtT(" ", "___")
	mtT("1", "__1__")
	mtT(" 1 ", "__1__")
	mtT("123 - abc", "__123 abc__")

	mtH := func(src string, wait []uint32) {
		tg := MakeTrgms(strings.NewReader(src))
		l := HashToList(tg)
		if !slice.Equal(l, wait) {
			fmt.Println(l)
			t.Fatal("Invalid trgms for: " + src)
		}
	}

	mtH("", []uint32{1582742881})
	mtH("1", []uint32{2770926408, 4099731774, 471817035})
	mtH("123 123", []uint32{2286445522, 2519459501, 3381053888, 3608265145, 397126932, 4099731774, 490120771, 530288421})

	eT := func(src1 string, src2 string, sim float64) {
		ret := Sim(strings.NewReader(src1), strings.NewReader(src2))
		if ret != sim {
			fmt.Println("expect sim ", sim, " real ", ret, " for ", src1, " and ", src2)
		}
	}

	eT("", "", 1)
	eT("1", "1", 1)
	eT("", "1", 0)
	eT("123", "123 123", 0.625)
}
