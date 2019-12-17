package romance

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

type isVowelFunc func(rune) bool

func VnvSuffix(word *snowballword.Word, f isVowelFunc, start int) int {
	for i := 1; i < len(word.RS[start:]); i++ {
		j := start + i
		if f(word.RS[j-1]) && !f(word.RS[j]) {
			return j + 1
		}
	}
	return len(word.RS)
}
