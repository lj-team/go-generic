package russian

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

// Step 3 is the removal of the derivational suffix.
//
func step3(word *snowballword.Word) bool {

	suffix, _ := word.RemoveFirstSuffixIn(word.R2start, "ост", "ость")
	if suffix != "" {
		return true
	}
	return false
}
