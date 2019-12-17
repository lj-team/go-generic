package bot

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/lj-team/go-generic/bot/dialog"
	"github.com/lj-team/go-generic/bot/history"
	"github.com/lj-team/go-generic/db/tindex"
	"github.com/lj-team/go-generic/text"
	"github.com/lj-team/go-generic/text/args"

	_ "github.com/lj-team/go-generic/db/tindex/engine/text"
)

type ANSWER_FUNC func(string, string) (string, bool)

type Bot struct {
	dlg       *dialog.Dialog
	history   *history.History
	index     tindex.Engine
	answer    ANSWER_FUNC
	threshold float64
}

type Config struct {
	RulesFile  string
	AnswerFunc ANSWER_FUNC
	History    string
	Threshold  float64
}

func New(cfg *Config) *Bot {

	b := &Bot{
		dlg:       dialog.New(),
		answer:    cfg.AnswerFunc,
		threshold: cfg.Threshold,
	}

	if cfg.History != "" {
		b.history, _ = history.New(cfg.History)
	}

	b.loadData(cfg.RulesFile)

	return b
}

func (b *Bot) Message(uniq string, msg string) (string, bool) {

	if msg == "" {
		return "", false
	}

	dlgID := b.dlg.Get(uniq)

	if b.history != nil {
		b.history.Write(dlgID, msg)
	}

	res := b.index.Search(prepareQuestion(msg), b.threshold)

	if len(res) > 0 {
		answer, ok := b.answer(dlgID, res[0].Value)
		if ok && b.history != nil {
			b.history.Write(dlgID, answer)
		}
		return answer, ok
	}

	return "", false
}

func (b *Bot) Close() {
	b.dlg.Close()
	b.index.Close()

	if b.history != nil {
		b.history.Close()
	}
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

func (b *Bot) load(rh io.Reader) {

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

func (b *Bot) loadData(filename string) error {

	b.index, _ = tindex.Open("text", "")

	rh, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer rh.Close()

	b.load(rh)

	return nil
}

func (b *Bot) AddRule(k, v string) {
	b.index.Add(k, v)
}
