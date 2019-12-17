package lang

import (
	"unicode"
)

const (
	LANG_RU int64 = 0x0001 // русские символы
	LANG_EN int64 = 0x0002 // английские символы
	LANG_UN int64 = 0x0004 // символы другого языка
	LANG_DI int64 = 0x0008 // цифры
)

var ruMap map[rune]bool
var enMap map[rune]bool

func init() {

	enMap = make(map[rune]bool)
	ruMap = make(map[rune]bool)

	for _, r := range "йцукенгшщзхъфывапролджэячсмитьбюё" {
		ruMap[r] = true
	}

	for _, r := range "qwertyuiopasdfghjklzxcvbnm" {
		enMap[r] = true
	}
}

func Lang(str string) int64 {

	mask := int64(0)

	for _, s := range str {

		if _, h := ruMap[s]; h {
			mask |= LANG_RU
			continue
		}

		if _, h := enMap[s]; h {
			mask |= LANG_EN
			continue
		}

		if unicode.IsLetter(s) {
			mask |= LANG_UN
			continue
		}

		if s >= '0' && s <= '9' {
			mask |= LANG_DI
			continue
		}
	}

	return mask
}
