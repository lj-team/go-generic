package stemmer

import (
	"io"
	"strconv"
	"strings"

	"github.com/lj-team/go-generic/text"
	"github.com/lj-team/go-generic/text/stemmer/dict"
	"github.com/lj-team/go-generic/text/stemmer/lang"
	"github.com/lj-team/go-generic/text/stemmer/lang/english"
	"github.com/lj-team/go-generic/text/stemmer/lang/russian"
)

var prefix map[string]bool

func init() {

	prefList := []string{"авиа", "авто", "агро", "аква", "анти", "арт", "арт-", "архи", "астро", "аудио", "аэро",
		"био", "бензо", "бого", "быстро",
		"вело", "веро", "видео", "вице-", "вне", "все", "втор", "высоко",
		"гео", "гетеро", "гидро", "гипер", "гос",
		"дерево", "дельта", "до", "долго", "древне",
		"евро",
		"загран", "звуко", "зоо", "инфра",
		"квадро", "квази", "кибер", "кино", "контр", "концепт", "кратко", "крио", "крипто",
		"лже",
		"макро", "мало", "мат", "мед", "мега", "медиа", "меж", "между", "место", "мета", "метео", "микро", "мини",
		"много", "моно", "мос", "мото", "мульти",
		"нано", "нарко", "нац", "не", "нейро", "нео", "нефте", "низко", "нитро", "ново",
		"пере", "пред", "пси", "пси-", "психо",
		"после", "пост", "право", "пред", "порно", "полу", "противо", "проф", "псевдо", "псевдо-",
		"радио", "ретро", "рос",
		"само", "сверх", "свето", "секс-", "сильно", "слабо", "смарт-", "соц", "социо", "спец", "спорт", "старо", "суб", "супер",
		"теле", "тепло", "термо", "терра", "тетра", "тех", "техно", "тихо", "тур",
		"ультра",
		"физико-", "фин", "фито", "фото",
		"центро",
		"экзо", "эко", "эконом-", "экс", "экс-", "экстра", "экстра-", "экшн-", "электро", "энерго", "этно"}

	prefix = map[string]bool{}

	for _, v := range prefList {
		prefix[v] = true
	}
}

// пробуем определить префикс для изестой основы
func tryPart(src string) (string, bool) {

	j := 0

	for i, _ := range src {

		if j > 0 {

			if _, h := prefix[src[:i]]; h {
				val := src[i:]

				if val, h = dict.Get(val); h {
					return src[:i] + val, true
				}
			}
		}

		if j > 6 {
			break
		}

		j++
	}

	return "", false
}

func TryExc(word string) (string, bool) {
	if wf, has := dict.Get(word); has {
		return wf, true
	}

	_, e := strconv.ParseInt(word, 10, 64)
	if e == nil {
		return word, true
	}

	if part, ok := tryPart(word); ok {
		return part, true
	}

	return "", false
}

func TryExcCombo(word string) (string, bool) {

	lst := strings.Split(word, "-")

	for i, v := range lst {

		if w, h := TryExc(v); h {
			lst[i] = w
		} else {
			return "", false
		}
	}

	return strings.Join(lst, "-"), true
}

// Stem a word
//
func Proc(word string) string {

	if wf, has := TryExc(word); has {
		return wf
	}

	for _, run := range word {
		if run == '-' {

			lst := strings.Split(word, "-")

			res := make([]string, 0, len(lst))

			for _, w := range lst {

				if wf, has := TryExc(w); has {
					res = append(res, wf)
					continue
				}

				msk := lang.Lang(w)

				if msk == lang.LANG_RU {
					res = append(res, russian.Stem(w))
				} else if msk == lang.LANG_EN {
					res = append(res, english.Stem(w))
				} else {
					res = append(res, w)
				}

			}

			return strings.Join(res, "-")
		}
	}

	msk := lang.Lang(word)

	if msk == lang.LANG_RU {
		return russian.Stem(word)
	} else if msk == lang.LANG_EN {

		return english.Stem(word)
	}

	return word
}

func TextToCode(str string) string {
	res := make([]string, 0, 4)

	for wrd := range Read(strings.NewReader(str)) {
		rec := Proc(wrd)
		res = append(res, rec)
	}

	return strings.Join(res, " ")
}

func Stream(input <-chan string) <-chan string {
	output := make(chan string, 4096)

	go func() {
		for w := range input {
			output <- Proc(w)
		}
		close(output)
	}()

	return output
}

func Read(in io.Reader) <-chan string {
	return text.Stream(in, Proc, 4096)
}

func LoadDict(filename string) error {
	return dict.Load(filename)
}
