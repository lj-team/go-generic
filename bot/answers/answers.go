package answers

import (
	"io"
	"io/ioutil"

	"github.com/lj-team/go-generic/rand/mt"
	"gopkg.in/yaml.v2"
)

type Answers struct {
	code2texts map[string][]string
}

func New(in io.Reader) (*Answers, error) {

	content, err := ioutil.ReadAll(in)

	hash := map[string][]string{}

	if err = yaml.Unmarshal(content, &hash); err != nil {
		return nil, err
	}

	return &Answers{
		code2texts: hash,
	}, nil
}

func (a *Answers) Get(code string) (string, bool) {

	if v, h := a.code2texts[code]; h {

		size := len(v)

		if size > 0 {

			r := mt.Next()
			if r < 0 {
				r = -r
			}

			return v[int(r%int64(size))], true

		}

	}

	return "", false
}
