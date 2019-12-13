package trgm

import (
	"hash/crc32"
	"io"
	"strings"

	"github.com/lj-team/go-generic/mapcnt"
	"github.com/lj-team/go-generic/text"
)

var (
	skipWrd map[string]bool
)

func init() {

	skipWrd = map[string]bool{
		".":  true,
		",":  true,
		"?":  true,
		"!":  true,
		";":  true,
		"…":  true,
		"(":  true,
		")":  true,
		"{":  true,
		"}":  true,
		"[":  true,
		"]":  true,
		"\"": true,
		"'":  true,
		"-":  true,
	}
}

func makeText(in io.Reader) string {

	maker := strings.Builder{}

	i := 0

	text.Proc(in, func(w string) {

		if skipWrd[w] {
			return
		}

		if i == 0 {
			maker.WriteRune('_')
			maker.WriteRune('_')
		} else {
			maker.WriteRune(' ')
		}

		i++

		maker.WriteString(w)
	})

	if i == 0 {
		maker.WriteString("___")
	} else {
		maker.WriteString("__")
	}

	return maker.String()
}

func MakeTrgms(in io.Reader) map[uint32]bool {

	res := make(map[uint32]bool)

	txt := makeText(in)
	i := 0

	w1 := ""
	w2 := ""
	w3 := ""

	for _, r := range txt {
		w1 = w2 + string(r)
		w2 = w3 + string(r)
		w3 = string(r)

		i++

		if i > 2 {
			res[crc32.ChecksumIEEE([]byte(w1))] = true
		}
	}

	return res

}

func UsedPart(txt io.Reader, source io.Reader) float64 {

	first := MakeTrgms(txt)
	second := MakeTrgms(source)

	return HashPart(first, second)
}

func HashPart(hash map[uint32]bool, source map[uint32]bool) float64 {
	has := 0
	size := len(hash)

	if size == 0 {
		return 0
	}

	for k := range hash {
		if _, h := source[k]; h {
			has++
		}
	}

	return float64(has) / float64(size)
}

func Sim(txt1 io.Reader, txt2 io.Reader) float64 {
	hash1 := MakeTrgms(txt1)
	hash2 := MakeTrgms(txt2)

	return HashSim(hash1, hash2)
}

func HashSim(hash1 map[uint32]bool, hash2 map[uint32]bool) float64 {
	has := 0

	size1 := len(hash1)
	if size1 == 0 {
		return 0
	}

	size2 := len(hash2)
	if size2 == 0 {
		return 0
	}

	for k := range hash1 {
		if _, h := hash2[k]; h {
			has++
		}
	}

	return float64(has) / float64(size1+size2-has)
}

func HashToList(hash map[uint32]bool) []uint32 {

	return mapcnt.List(hash, nil).([]uint32)

}
