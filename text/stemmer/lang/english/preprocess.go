package english

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

func preprocess(word *snowballword.Word) {

	// Clean up apostrophes
	normalizeApostrophes(word)
	trimLeftApostrophes(word)

	// Capitalize Y's that are not behaving
	// as vowels.
	capitalizeYs(word)

	// Find the two regions, R1 & R2
	r1start, r2start := r1r2(word)
	word.R1start = r1start
	word.R2start = r2start
}
