package text

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type StrFunc func(string)
type StrConvFunc func(string) string

var (
	skipWrd map[string]bool
)

func init() {

	skipWrd = map[string]bool{
		".":  true,
		",":  true,
		"?":  true,
		"!":  true,
		":":  true,
		";":  true,
		"…":  true,
		"(":  true,
		")":  true,
		"{":  true,
		"}":  true,
		"[":  true,
		"]":  true,
		"\"": true,
		"'":  true,
		"-":  true,
	}
}

func Proc(in io.Reader, fn StrFunc) {

	reader := bufio.NewReader(in)
	state := 0
	builder := strings.Builder{}
	prev := rune(0)

	sendSep := func(run rune) {
		fn(string(run))
	}

	sendEos := func(run rune) {
		fn(string(run))
	}

	finishWord := func() {
		fn(builder.String())
	}

	proc0 := func(run rune) {
		if unicode.IsLetter(run) || run == '_' {
			builder.Reset()
			builder.Grow(32)
			builder.WriteRune(run)
			state = 1
		} else if unicode.IsDigit(run) {
			builder.Reset()
			builder.Grow(16)
			builder.WriteRune(run)
			state = 2
		} else if s := isSep(run); s != 0 {
			sendSep(s)
		} else if s := isEos(run); s != 0 {
			sendEos(s)
		}
	}

	for {

		run, _, err := reader.ReadRune()
		if err != nil && run == 0 {
			break
		}

		run = unicode.ToLower(run)
		if run == 'ё' {
			run = 'е'
		}

		switch state {

		// начальное состояние
		case 0:

			proc0(run)

			// загрузка слова
		case 1:

			if unicode.IsLetter(run) || unicode.IsDigit(run) || run == '_' {
				builder.WriteRune(run)
			} else if run == ':' && isProto(builder.String()) {
				state = 3
			} else if run == '-' || run == '.' {
				prev = run
				state = 5
			} else if s := isSep(run); s != 0 {
				finishWord()
				sendSep(s)
				state = 0
			} else if s := isEos(run); s != 0 {
				finishWord()
				sendEos(s)
				state = 0
			} else {
				finishWord()
				state = 0
			}

		// загрузка числа
		case 2:
			if unicode.IsDigit(run) {
				builder.WriteRune(run)
			} else if unicode.IsLetter(run) || run == '_' {
				builder.WriteRune(run)
				state = 1
			} else if run == '-' {
				state = 6
			} else if s := isSep(run); s != 0 {
				finishWord()
				sendSep(s)
				state = 0
			} else if s := isEos(run); s != 0 {
				finishWord()
				sendEos(s)
				state = 0
			} else {
				finishWord()
				state = 0
			}

			// начало загрузки урла
		case 3:

			if run == '/' {
				builder.WriteByte(':')
				builder.WriteByte('/')
				state = 4
			} else {

				finishWord()
				sendSep(':')

				state = 0

				proc0(run)
			}

		// чтение урла до конца
		case 4:

			if !IsSpace(run) {
				builder.WriteRune(run)
			} else {
				finishWord()
				state = 0
			}

		// проверка того, что идет в слове после дефиса/точки, если букцы или цифры, то слово продолжается
		case 5:
			if unicode.IsLetter(run) || unicode.IsDigit(run) || run == '_' {
				builder.WriteRune(prev)
				builder.WriteRune(run)
				state = 1
			} else {
				finishWord()
				if prev == '-' {
					sendSep(prev)
				} else {
					sendEos(prev)
				}

				state = 0

				proc0(run)
			}

			//диапазон 1941-1945
		case 6:
			if unicode.IsDigit(run) {
				builder.WriteByte('-')
				builder.WriteRune(run)
				state = 7
			} else if unicode.IsLetter(run) || run == '_' {
				builder.WriteByte('-')
				builder.WriteRune(run)
				state = 1
			} else {
				finishWord()
				sendSep('-')
				state = 0
				proc0(run)
			}

		case 7:

			if unicode.IsDigit(run) {
				builder.WriteRune(run)
			} else if unicode.IsLetter(run) || run == '_' {
				builder.WriteRune(run)
				state = 1
			} else if run == '-' {
				state = 8
			} else if s := isSep(run); s != 0 {
				finishWord()
				sendSep(s)
				state = 0
			} else if s := isEos(run); s != 0 {
				finishWord()
				sendEos(s)
				state = 0
			} else {
				finishWord()
				state = 0
			}

			//диапазон 28-05-1985
		case 8:
			if unicode.IsDigit(run) {
				builder.WriteByte('-')
				builder.WriteRune(run)
				state = 9
			} else if unicode.IsLetter(run) || run == '_' {
				builder.WriteByte('-')
				builder.WriteRune(run)
				state = 1
			} else {
				finishWord()
				sendSep('-')
				state = 0
				proc0(run)
			}

		case 9:

			if unicode.IsDigit(run) {
				builder.WriteRune(run)
			} else if unicode.IsLetter(run) || run == '_' {
				builder.WriteRune(run)
				state = 1
			} else if run == '-' {
				prev = '-'
				state = 5
			} else if s := isSep(run); s != 0 {
				finishWord()
				sendSep(s)
				state = 0
			} else if s := isEos(run); s != 0 {
				finishWord()
				sendEos(s)
				state = 0
			} else {
				finishWord()
				state = 0
			}

		}
	}

	switch state {

	// последнее слово
	case 1:

		finishWord()

	// последнее число
	case 2:

		finishWord()

	// http: or https:
	case 3:

		finishWord()
		sendSep(':')

	// загрузка урла
	case 4:

		finishWord()

	// слово закончили дефисом или точкой
	case 5:

		finishWord()
		if prev == '-' {
			sendSep(prev)
		} else {
			sendEos(prev)
		}

	case 6:

		finishWord()
		sendSep('-')

	// 1941-1945
	case 7:

		finishWord()

	case 8:

		finishWord()
		sendSep('-')

	// 28-05-1985
	case 9:

		finishWord()

	}
}

func Stream(in io.Reader, fn StrConvFunc, chanSize int) <-chan string {

	out := make(chan string, chanSize)

	go func() {
		defer close(out)

		Proc(in, func(str string) {
			if val := fn(str); val != "" {
				out <- val
			}
		})

	}()

	return out
}

func WordsOnly(in io.Reader) string {

	i := 0
	maker := strings.Builder{}

	Proc(in, func(w string) {

		if skipWrd[w] {
			return
		}

		if i > 0 {
			maker.WriteRune(' ')
		}
		maker.WriteString(w)

		i++
	})

	return maker.String()
}
