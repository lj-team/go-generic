package history

import (
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/lj-team/go-generic/time/strftime"
)

type History struct {
	level     int
	template  string
	save      int64
	period    int64
	input     chan string
	fh        *os.File
	filename  string
	lastCheck int64
	finish    chan bool
}

var replace map[rune]rune = map[rune]rune{
	'ё': 'е',
	'Ё': 'Е',
}

func replaceRune(r rune) rune {

	if v, h := replace[r]; h {
		return v
	}

	return r
}

func textPreproc(msg string) string {

	str := strings.Builder{}
	src := strings.NewReader(msg)
	mod := 0

	for {

		ch, _, e := src.ReadRune()
		if e != nil && ch == 0 {
			break
		}

		switch mod {
		case 0:

			if !unicode.IsSpace(ch) {
				str.WriteRune(replaceRune(ch))
				mod = 1
			}

		case 1:

			if unicode.IsSpace(ch) {
				mod = 2
			} else {
				str.WriteRune(replaceRune(ch))
			}

		case 2:

			if !unicode.IsSpace(ch) {
				str.WriteRune(' ')
				str.WriteRune(replaceRune(ch))
				mod = 1
			}
		}
	}

	return str.String()
}

func (h *History) Write(dlgID string, msg string) {

	if h == nil {
		return
	}

	text := dlgID + " " + textPreproc(msg)

	h.input <- text
}

func New(path string) (*History, error) {

	l := &History{
		template:  path,
		period:    86400,
		save:      31,
		filename:  strftime.Format(path, time.Now()),
		lastCheck: time.Now().Unix(),
		input:     make(chan string, 1024),
		finish:    make(chan bool),
	}

	var err error

	if l.fh, err = os.OpenFile(l.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err != nil {
		return nil, err
	}

	if l.save > 0 {
		rm_name := strftime.Format(l.template, time.Unix(l.lastCheck-int64(l.save*l.period), 0))
		os.Remove(rm_name)
	}

	go l.writer()

	return l, nil
}

func (l *History) rotate() {

	if l.lastCheck+60 > time.Now().Unix() {
		return
	}

	l.lastCheck = time.Now().Unix()
	new_name := strftime.Format(l.template, time.Now())

	if new_name != l.filename {
		l.fh.Close()

		var err error
		l.filename = new_name

		if l.fh, err = os.OpenFile(l.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err != nil {
			panic(err)
		}

		if l.save > 0 {
			rm_name := strftime.Format(l.template, time.Unix(l.lastCheck-int64(l.save*l.period), 0))
			os.Remove(rm_name)
		}
	}
}

func (l *History) writer() {
	for {
		select {

		case str, ok := <-l.input:

			l.rotate()

			if str != "" {
				l.fh.WriteString(str + "\n")
				l.fh.Sync()
			}

			if !ok {
				close(l.finish)
				return
			}

		case <-time.After(time.Minute):
			l.rotate()
		}
	}
}

func (l *History) Close() {
	close(l.input)
	<-l.finish
}
