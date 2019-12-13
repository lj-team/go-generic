package text

import (
	"strings"
)

func isSep(run rune) rune {
	if r, h := sep[run]; h {
		return r
	}
	return 0
}

func isEos(run rune) rune {
	if r, h := eos[run]; h {
		return r
	}
	return 0
}

func IsSpace(run rune) bool {
	return run == ' ' || run == '\t' || run == '\n' || run == '\r'
}

func isProto(str string) bool {
	_, h := proto[str]
	return h
}

func Truncate(text string, limit int) string {

	builder := strings.Builder{}
	i := 0

	for _, rune := range text {
		builder.WriteRune(rune)
		i++

		if i >= limit {
			builder.WriteRune('…')
			break
		}
	}
	return builder.String()
}

func Quoted(text string) string {

	var sb strings.Builder

	sb.WriteRune('«')
	sb.WriteString(text)
	sb.WriteRune('»')

	return sb.String()
}

func Reverse(text string) string {

	var sb strings.Builder

	list := []rune(text)

	for i := len(list) - 1; i > -1; i-- {
		sb.WriteRune(list[i])
	}

	return sb.String()
}

func CmpSlices(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}

// add value if not exists
func AddNew(list []string, str string) []string {

	for _, v := range list {
		if v == str {
			return list
		}
	}

	list = append(list, str)
	return list
}
