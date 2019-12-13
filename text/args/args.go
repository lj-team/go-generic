package args

import (
	"errors"
	"io"
	"strings"
)

var (
	mapper    map[byte]byte
	mapperRev map[rune]string
)

func init() {

	mapper = map[byte]byte{
		'\\': '\\',
		'\'': '\'',
		'"':  '"',
		's':  ' ',
		't':  '\t',
		'r':  '\r',
		'n':  '\n',
		' ':  ' ',
	}

	mapperRev = map[rune]string{
		' ':  " ",
		'\t': "\\t",
		'\n': "\\n",
		'\r': "\\r",
		'"':  "\\\"",
		'\'': "\\'",
		'\\': "\\\\",
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func parse(rd io.ByteReader) ([]string, error) {

	mode := 0
	bkmode := 0
	cur := []byte{}
	var res []string

	for {

		c, err := rd.ReadByte()
		if c == 0 && err != nil {
			break
		}

		switch mode {
		case 0:
			if isSpace(c) {
				continue
			}

			if c == '"' {
				mode = 1
				cur = []byte{}
				continue
			}

			cur = []byte{c}
			mode = 2

		case 1:

			if c == '"' {
				mode = 4
				res = append(res, string(cur))
				continue
			}

			if c == '\\' {
				mode = 3
				bkmode = 1
				continue
			}

			cur = append(cur, c)

		case 2:

			if isSpace(c) {
				res = append(res, string(cur))
				mode = 0
				continue
			}

			if c == '\\' {
				mode = 3
				bkmode = 2
				continue
			}

			cur = append(cur, c)

		case 3:

			nc, ok := mapper[c]

			if !ok {
				return nil, errors.New("inavlid format")
			}

			cur = append(cur, nc)
			mode = bkmode

		case 4:
			if isSpace(c) {
				mode = 0
				continue
			}

			return nil, errors.New("inavlid format")
		}
	}

	if mode == 1 || mode == 3 {
		return nil, errors.New("inavlid format")
	}

	if mode == 2 {
		res = append(res, string(cur))
	}

	return res, nil
}

func Read(input io.ByteReader) []string {

	tokens, err := parse(input)

	if err != nil {
		return nil
	}

	return tokens
}

func Parse(str string) []string {

	tokens, err := parse(strings.NewReader(str))

	if err != nil {
		return nil
	}

	return tokens
}

func encode(str string) string {

	if str == "" {
		return `""`
	}

	found := false

	for _, r := range str {

		if _, h := mapperRev[r]; h {
			found = true
			break
		}

	}

	if !found {
		return str
	}

	builder := strings.Builder{}

	builder.WriteRune('"')

	for _, r := range str {

		if nv, h := mapperRev[r]; h {
			builder.WriteString(nv)
		} else {
			builder.WriteRune(r)
		}

	}

	builder.WriteRune('"')

	return builder.String()
}

func ToString(args ...string) string {

	builder := strings.Builder{}

	for i, str := range args {

		if i > 0 {
			builder.WriteRune(' ')
		}

		builder.WriteString(encode(str))
	}

	return builder.String()
}
