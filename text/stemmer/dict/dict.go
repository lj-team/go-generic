package dict

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/lj-team/go-generic/resource"
	"github.com/lj-team/go-generic/text/args"

	_ "github.com/lj-team/go-generic/text/stemmer/data"
)

var (
	dict map[string]string
)

func init() {
	dict = map[string]string{}

	data := resource.Get("github.com/lj-team/go-generic/text/stemmer/data/stemmer.txt")

	if data != nil {
		load(bytes.NewReader(data.Data))
	}
}

func Get(wrd string) (string, bool) {
	v, h := dict[wrd]

	if !h {
		return wrd, false
	}

	return v, h
}

func Push(wrd string, val string) {
	dict[wrd] = val
}

func Load(filename string) error {

	rf, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer rf.Close()

	load(rf)

	return nil
}

func load(rh io.Reader) {

	br := bufio.NewReader(rh)

	data := map[string]string{}

	for {

		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			break
		}

		lst := args.Parse(str)

		switch len(lst) {

		case 1:

			data[lst[0]] = lst[0]

		case 2:

			data[lst[0]] = lst[1]
		}

	}

	dict = data
}
