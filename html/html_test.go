package html

import (
	"strings"
	"testing"
)

func TestTextHtmlProcessString(t *testing.T) {
	h := New(DefaultOpts())

	if h == nil {
		t.Fatal("NewHtmlParser() erro")
	}

	src := `<html><head><title>При<WBR/>вет ми<WbR />р</title><style>a{color:#CCC;}</style></head><body><h1>He<wbr>l&shy;l&#173;o<del>1</del></h1><script><!--
	alert('hello')
--></script><iframe src="xxx"></iframe><strike>123</strike><a href='/url1'>url1</a><a href='/url2'><s>1</s>url2</a><img src="" title="img title"/><img src="" title="привет"/><img src="" title="dsc123"/><img src="" title="Смайлотрон - смайлики для ЖЖ!"/></body></html>`

	res := h.Process(strings.NewReader(src))

	if res != " Привет мир Hello url1 url2 .\nimg title .\n .\nпривет .\n" {
		t.Fatal("Wrong plain text")
	}

	if len(h.iframes) != 1 || !h.iframes["xxx"] {
		t.Fatal("Iframe src not found")
	}

	has := false

	for v := range h.Iframes() {
		if v != "xxx" {
			t.Fatal("bad iframes")
		}

		has = true
	}

	if !has {
		t.Fatal("empty chan")
	}

	opts := DefaultOpts()
	opts.SocialCleanText = true

	h = New(opts)

	src = `<body>Follow me: <a>Twitter</a>/<a>Facebook</a>/<a>YouTube</a>!</body>`
	res = h.Process(strings.NewReader(src))

	if res != ` Follow me:  / / !` {
		t.Fatal("Process failed")
	}
}
