package kv

import (
	"strings"
)

func encode(str string) string {

	if str == "" {
		return str
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

	for _, r := range str {

		if nv, h := mapperRev[r]; h {
			builder.WriteString(nv)
		} else {
			builder.WriteRune(r)
		}

	}

	return builder.String()
}
