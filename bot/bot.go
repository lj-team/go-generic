package bot

import (
	"bufio"
	"io"
	"strings"

	"github.com/lj-team/go-generic/bot/answers"
	"github.com/lj-team/go-generic/bot/dialog"
	"github.com/lj-team/go-generic/cache"
	"github.com/lj-team/go-generic/db/tindex"
	"github.com/lj-team/go-generic/log"
	"github.com/lj-team/go-generic/text"
	"github.com/lj-team/go-generic/text/args"

	_ "github.com/lj-team/go-generic/db/tindex/engine/text"
)

type AnswerFunc func(string, string) (string, bool)

type Bot struct {
	name   string
	index  tindex.Engine
	answer AnswerFunc
	sens   float64
}

func New(botName string, schema io.Reader, ans io.Reader, sens float64) (*Bot, error) {

	af, err := makeAnswerFunc(ans)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		name:   botName,
		answer: af,
		sens:   sens,
	}

	b.index, _ = tindex.Open("text", "")

	b.loadSchema(schema)

	return b, nil
}

func (b *Bot) Message(uniq string, msg string) (string, bool) {

	if msg == "" {
		return "", false
	}

	dlgID := dialog.Get(uniq)

	log.Infof("bot=%s dlg=%s in=%s", b.name, dlgID, msg)

	res := b.index.Search(prepareQuestion(msg), b.sens)

	if len(res) > 0 {
		answer, ok := b.answer(dlgID, res[0].Value)
		if ok {
			log.Infof("bot=%s dlg=%s out=%s", b.name, dlgID, answer)
		}
		return answer, ok
	}

	return "", false
}

func (b *Bot) Close() {
	b.index.Close()
}

func prepareQuestion(txt string) string {

	bldr := strings.Builder{}
	mode := true

	seps := map[string]bool{
		".":  true,
		",":  true,
		"!":  true,
		"?":  true,
		":":  true,
		";":  true,
		"(":  true,
		")":  true,
		"\"": true,
		"'":  true,
		"«":  true,
		"»":  true,
		"-":  true,
	}

	text.Proc(strings.NewReader(txt), func(item string) {

		if _, h := seps[item]; h {

			bldr.WriteString(item)
			mode = true

		} else {

			if mode {
				mode = false
			} else {
				bldr.WriteRune(' ')
			}
			bldr.WriteString(item)
		}

	})

	return bldr.String()
}

func (b *Bot) loadSchema(rh io.Reader) {

	br := bufio.NewReader(rh)

	for {

		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			break
		}

		tok := args.Parse(str)
		if len(tok) != 2 {
			continue
		}

		b.index.Add(prepareQuestion(tok[0]), tok[1])
	}

}

func (b *Bot) AddRule(k, v string) {
	b.index.Add(k, v)
}

func makeAnswerFunc(in io.Reader) (AnswerFunc, error) {

	a, err := answers.New(in)
	if err != nil {
		return nil, err
	}

	ch := cache.New("size=10240 nodes=16 ttl=86400")

	return func(dlgID string, msg string) (string, bool) {
		for i := 0; i < 5; i++ {

			txt, ok := a.Get(msg)

			if ok {

				key := dlgID + ":" + txt

				if ch.Get(key) != nil {
					continue
				}

				ch.Set(key, "1")
				return txt, ok
			}
		}

		return "", false
	}, nil
}
