package english

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

// Applies transformations necessary after
// a word has been completely processed.
//
func postprocess(word *snowballword.Word) {

	uncapitalizeYs(word)
}
