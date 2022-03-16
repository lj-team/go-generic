package html

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/lj-team/go-generic/text/args"
	"github.com/lj-team/go-generic/text/lang"
	ht "golang.org/x/net/html"
)

var (
	scClean map[string]bool
)

func init() {
	scClean = make(map[string]bool)

	for _, w := range []string{"facebook", "instagram", "livejournal", "ok", "tripadviser", "twitter", "youtube",
		"вк", "дзен", "живой журнал", "инстаграм", " одноклассники", "ок", "твиттер", "телеграм-канал",
		"ютюб", "яндекс дзен", "яндекс.дзен"} {
		scClean[w] = true
	}
}

type HTML struct {
	links   map[string]bool
	iframes map[string]bool
	images  map[string]bool
	divs    []map[string]string

	Erase map[string]bool

	tagStack   []string
	eraseStack []string

	buf *bytes.Buffer

	opts *Opts
}

type Opts struct {
	ReadImgAlt      bool
	OnlyCleanHtml   bool
	LinkStat        bool
	EraseTags       []string
	SocialCleanText bool // remove `follow me` social networks names from result text
}

func DefaultOpts() *Opts {

	return &Opts{
		ReadImgAlt: true,
		LinkStat:   true,
		EraseTags:  []string{"form", "iframe", "link", "meta", "noscript", "option", "script", "select", "style", "s", "strike", "del"},
	}

}

func New(o *Opts) *HTML {

	res := &HTML{}

	if o == nil {
		o = DefaultOpts()
	}

	res.Erase = map[string]bool{}

	for _, v := range o.EraseTags {
		res.Erase[v] = true
	}

	res.tagStack = make([]string, 0)
	res.eraseStack = make([]string, 0)

	res.buf = new(bytes.Buffer)

	res.opts = o

	return res
}

func (h *HTML) Process(r io.Reader) string {

	data, _ := ioutil.ReadAll(r)

	re, err := regexp.Compile("(?i)<wbr\\s*/?>|&shy;|&#173;|</wbr>")
	if err == nil {
		data = re.ReplaceAll(data, []byte(""))
	}

	h.buf.Reset()

	h.tagStack = h.tagStack[:0]
	h.eraseStack = h.eraseStack[:0]

	h.links = make(map[string]bool)
	h.iframes = make(map[string]bool)
	h.images = make(map[string]bool)

	parser := ht.NewTokenizer(bytes.NewReader(data))

	for {
		tt := parser.Next()

		switch {
		case tt == ht.ErrorToken:
			return h.buf.String()

		case tt == ht.StartTagToken:
			t := parser.Token()
			h.onStartTag(&t, string(parser.Raw()))

		case tt == ht.EndTagToken:
			t := parser.Token()
			h.onCloseTag(&t)

		case tt == ht.SelfClosingTagToken:
			t := parser.Token()
			h.onStartTag(&t, string(parser.Raw()))
			h.onCloseTag(&t)

		case tt == ht.TextToken:
			h.onText(string(parser.Text()), string(parser.Raw()))

		}
	}

	return h.buf.String()
}

func (h *HTML) onStartTag(t *ht.Token, raw string) {

	if len(h.eraseStack) > 0 {
		h.eraseStack = append(h.eraseStack, t.Data)
		return
	}

	_, found := h.Erase[t.Data]

	if found {
		if t.Data != "meta" && t.Data != "link" {
			h.eraseStack = append(h.eraseStack, t.Data)
		}

		if h.opts.LinkStat {

			switch t.Data {
			case "img":

				for _, a := range t.Attr {
					if a.Key == "src" {
						h.images[a.Val] = true
					}
				}

			case "a":

				for _, a := range t.Attr {
					if a.Key == "href" {
						h.iframes[a.Val] = true
					}
				}

			case "iframe":

				for _, a := range t.Attr {
					if a.Key == "src" {
						h.iframes[a.Val] = true
					}
				}

			}

		}

		return
	}

	h.tagStack = append(h.tagStack, t.Data)

	switch t.Data {

	case "a":

		for _, a := range t.Attr {
			if a.Key == "href" && h.opts.LinkStat {
				h.links[a.Val] = true
			}
		}

	case "iframe":

		for _, a := range t.Attr {
			if a.Key == "src" && h.opts.LinkStat {
				h.iframes[a.Val] = true
			}
		}

	case "div":

		for _, a := range t.Attr {
			tmp := make(map[string]string)
			tmp[a.Key] = a.Val
			h.divs = append(h.divs, tmp)
		}

	case "img":

		if h.opts.ReadImgAlt && !h.opts.OnlyCleanHtml {

			title := ""

			for _, a := range t.Attr {

				switch a.Key {
				case "title", "alt":

					if !skipAlt(a.Val) && title != a.Val {

						title = a.Val

						h.buf.WriteString(" .\n")
						h.buf.WriteString(title)
						h.buf.WriteString(" .\n")

						title = a.Val
					}

				case "src":

					if h.opts.LinkStat {
						h.images[a.Val] = true
					}

				}
			}
		}
	}

	if h.opts.OnlyCleanHtml {
		h.buf.WriteString(raw)
	}
}

func skipAlt(txt string) bool {

	if txt == "" {
		return true
	}

	list := args.Parse(txt)

	if len(list) == 0 {
		return true
	}

	if strings.Index(txt, "смайлики для ЖЖ") >= 0 {
		return true
	}

	if len(list) == 1 {

		if !lang.IsRusWord(list[0]) && !lang.IsEngWord(list[0]) {
			return true
		}

	} else {

		if len(list) > 2 && !lang.IsRus(txt) {
			return true
		}

	}

	return false
}

func (h *HTML) onCloseTag(t *ht.Token) {

	elen := len(h.eraseStack)

	if elen > 0 {

		if h.eraseStack[elen-1] == t.Data {
			h.eraseStack = h.eraseStack[0 : elen-1]
			return
		}

		pos := -1

		for i, cur := range h.eraseStack {
			if cur == t.Data {
				pos = i
			}
		}

		if pos > -1 {
			h.eraseStack = h.eraseStack[0:pos]
		}

		return
	}

	elen = len(h.tagStack)

	if len(h.tagStack) > 0 {

		if h.tagStack[elen-1] == t.Data {

			if h.opts.OnlyCleanHtml {
				h.buf.WriteString("</")
				h.buf.WriteString(t.Data)
				h.buf.WriteRune('>')
			}

			h.tagStack = h.tagStack[0 : elen-1]
			return
		}

		_, found := h.Erase[t.Data]
		if found {
			return
		}

		pos := -1

		for i, cur := range h.tagStack {
			if cur == t.Data {
				pos = i
			}
		}

		if pos > -1 {

			for i := elen - 1; i >= pos; i-- {

				h.buf.WriteString("</")
				h.buf.WriteString(h.tagStack[i])
				h.buf.WriteRune('>')

			}

			h.tagStack = h.tagStack[0:pos]
		}
	}
}

func (h *HTML) onText(str string, raw string) {

	if len(h.eraseStack) == 0 {

		if h.opts.OnlyCleanHtml {
			h.buf.WriteString(raw)
		} else {

			if h.opts.SocialCleanText {
				if size := len(h.tagStack); size > 0 {
					if h.tagStack[size-1] == "a" {
						lstr := strings.TrimSpace(strings.ToLower(str))
						if scClean[lstr] {
							return
						}
					}
				}
			}

			h.buf.WriteRune(' ')
			h.buf.WriteString(str)
		}

	}
}

func (h *HTML) iterate(hash map[string]bool) <-chan string {
	out := make(chan string)

	go func() {

		defer close(out)

		for v := range hash {

			out <- v
		}

	}()

	return out
}

func (h *HTML) iterateArr(arr []map[string]string) <-chan map[string]string {
	out := make(chan map[string]string)

	go func() {

		defer close(out)

		for _, v := range arr {

			out <- v
		}

	}()

	return out
}

func (h *HTML) Images() <-chan string {

	return h.iterate(h.images)
}

func (h *HTML) Links() <-chan string {

	return h.iterate(h.links)
}

func (h *HTML) Iframes() <-chan string {

	return h.iterate(h.iframes)
}

func (h *HTML) Divs() <-chan map[string]string {

	return h.iterateArr(h.divs)
}
