package lang

import (
	"strings"
	"unicode"
)

var (
	ruMap    map[rune]bool
	cyrMap   map[rune]bool
	enMap    map[rune]bool
	otherMap map[rune]bool
)

func init() {

	enMap = make(map[rune]bool)
	ruMap = make(map[rune]bool)
	otherMap = make(map[rune]bool)
	cyrMap = make(map[rune]bool)

	for _, r := range "йцукенгшщзхъфывапролджэячсмитьбюё" {
		ruMap[r] = true
	}

	for _, r := range "qwertyuiopasdfghjklzxcvbnm" {
		enMap[r] = true
	}

	for _, r := range "іўєაბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰәҙөғҫҡңһѣ" {
		otherMap[unicode.ToLower(r)] = true
	}

	for _, r := range "йцукенгшщзхъфывапролджэячсмитьбюёіўєһѣ" {
		cyrMap[r] = true
	}
}

func IsCyr(str string) bool {

	str = strings.ToLower(str)

	cyrCnt := 0
	nonCyr := 0

	for _, r := range str {
		if _, h := cyrMap[r]; h {
			cyrCnt++
		} else if unicode.IsLetter(r) {
			nonCyr++
		}
	}

	return cyrCnt > nonCyr
}

func IsOtherCyr(str string) bool {

	str = strings.ToLower(str)

	cntByOther := 0

	for _, s := range str {
		if _, h := otherMap[s]; h {
			cntByOther++
		}
	}

	return cntByOther > 5
}

func IsEng(str string) bool {

	str = strings.ToLower(str)

	engCnt := 0
	othCnt := 0

	for _, s := range str {

		if _, h := enMap[s]; h {
			engCnt++
			continue
		}

		if unicode.IsLetter(s) {
			othCnt++
		}

	}

	return engCnt > othCnt*2
}

func IsRus(str string) bool {

	str = strings.ToLower(str)

	cntRu := 0
	cntEn := 0
	cntByOther := 0
	cntSim := 0

	for _, s := range str {
		if _, h := ruMap[s]; h {
			cntRu++
			continue
		}

		if _, h := enMap[s]; h {
			cntEn++
			continue
		}

		if _, h := otherMap[s]; h {
			cntByOther++
		}

		if unicode.IsLetter(s) {
			cntSim++
		}
	}

	if cntByOther > 6 {
		return false
	}

	return cntRu >= cntEn && cntRu > cntSim*2
}

func IsRusWord(w string) bool {

	w = strings.ToLower(w)

	for _, r := range w {

		if _, h := ruMap[r]; !h && r != '-' {
			return false
		}
	}

	return true
}

func IsEngWord(w string) bool {

	w = strings.ToLower(w)

	for _, r := range w {

		if _, h := enMap[r]; !h && r != '-' {
			return false
		}
	}

	return true
}
