package kv

import (
	"errors"
)

type KV map[string]string

var (
	mapper    map[rune]rune
	mapperRev map[rune]string

	ErrParse            error
	ErrKeyAlreadyExists error
)

func init() {

	ErrParse = errors.New("parse error")
	ErrKeyAlreadyExists = errors.New("key already exists")

	mapper = map[rune]rune{
		'\\': '\\',
		'\'': '\'',
		'"':  '"',
		's':  ' ',
		't':  '\t',
		'r':  '\r',
		'n':  '\n',
		' ':  ' ',
		'=':  '=',
	}

	mapperRev = map[rune]string{
		' ':  "\\s",
		'\t': "\\t",
		'\n': "\\n",
		'\r': "\\r",
		'"':  "\\\"",
		'\'': "\\'",
		'\\': "\\\\",
		'=':  "\\=",
	}
}
