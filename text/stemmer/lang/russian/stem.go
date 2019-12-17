package russian

import (
	"strings"

	snowballword "github.com/lj-team/go-generic/text/stemmer/word"
)

func Stem(word string) string {

	word = strings.ToLower(strings.TrimSpace(word))
	w := snowballword.New(word)

	if len(w.RS) <= 2 || isStopWord(word) {
		return word
	}

	preprocess(w)
	step1(w)
	step2(w)
	step3(w)
	step4(w)
	return w.String()

}
