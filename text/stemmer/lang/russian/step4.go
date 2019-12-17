package russian

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

func step4(word *snowballword.Word) bool {

	// (1) Undouble "н", or, 2) if the word ends with a SUPERLATIVE ending,
	// (remove it and undouble н n), or 3) if the word ends ь (') (soft sign)
	// remove it.

	// Undouble "н"
	if word.HasSuffixRunes([]rune("нн")) {
		word.RemoveLastNRunes(1)
		return true
	}

	// Remove superlative endings
	suffix, _ := word.RemoveFirstSuffix("ейше", "ейш")
	if suffix != "" {
		// Undouble "н"
		if word.HasSuffixRunes([]rune("нн")) {
			word.RemoveLastNRunes(1)
		}
		return true
	}

	// Remove soft sign
	if rsLen := len(word.RS); rsLen > 0 && word.RS[rsLen-1] == 'ь' {
		word.RemoveLastNRunes(1)
		return true
	}
	return false
}
