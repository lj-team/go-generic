package russian

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

func step2(word *snowballword.Word) bool {
	suffix, _ := word.RemoveFirstSuffixIn(word.RVstart, "и")
	if suffix != "" {
		return true
	}
	return false
}
