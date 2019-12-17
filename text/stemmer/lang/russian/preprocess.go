package russian

import (
	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

func preprocess(word *snowballword.Word) {

	r1start, r2start, rvstart := findRegions(word)
	word.R1start = r1start
	word.R2start = r2start
	word.RVstart = rvstart

}
