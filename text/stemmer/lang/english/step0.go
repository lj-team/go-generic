package english

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

// Step 0 is to strip off apostrophes and "s".
//
func step0(w *snowballword.Word) bool {
	suffix, suffixRunes := w.FirstSuffix("'s'", "'s", "'")
	if suffix == "" {
		return false
	}
	w.RemoveLastNRunes(len(suffixRunes))
	return true
}
