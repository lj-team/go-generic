package text

import (
	"strings"
	"testing"

	"github.com/lj-team/go-generic/slice"
)

func TestTokens(t *testing.T) {

	tF := func(txt string, wait []string) {

		var res []string

		for w := range Stream(strings.NewReader(txt), func(w string) string { return w }, 128) {
			res = append(res, w)
		}

		if !slice.Equal(res, wait) {
			t.Fatal("test failed for: " + txt)
		}
	}

	tF("доброе утро", []string{"доброе", "утро"})
	tF("привет, мир!", []string{"привет", ",", "мир", "!"})
	tF("2 стакана", []string{"2", "стакана"})
	tF("9A - это мой класс", []string{"9a", "-", "это", "мой", "класс"})
	tF("12 негретят", []string{"12", "негретят"})
	tF("в 12", []string{"в", "12"})
	tF("перейди на https://yandex.ru и", []string{"перейди", "на", "https://yandex.ru", "и"})
	tF("перейди на https://yandex.ru", []string{"перейди", "на", "https://yandex.ru"})
	tF("хороший день.", []string{"хороший", "день", "."})
	tF("наука- это", []string{"наука", "-", "это"})
	tF("домен mail.ru", []string{"домен", "mail.ru"})
	tF("в 1941-1945 годы", []string{"в", "1941-1945", "годы"})
	tF("дата 28-05-1985", []string{"дата", "28-05-1985"})
	tF("номер 28-3ГЛ", []string{"номер", "28-3гл"})
	tF("номер 28-38ГЛ", []string{"номер", "28-38гл"})
	tF("номер 28-38-44-12", []string{"номер", "28-38-44-12"})
	tF("xx век", []string{"xx", "век"})
	tF("это только игра", []string{"это", "только", "игра"})
	tF("Привет, мир!", []string{"привет", ",", "мир", "!"})
	tF("Перешел на домен belfinor.ru", []string{"перешел", "на", "домен", "belfinor.ru"})
	tF("Привет _, как дела?", []string{"привет", "_", ",", "как", "дела", "?"})
	tF("Напиши _b_c_d, точка!", []string{"напиши", "_b_c_d", ",", "точка", "!"})
	tF("Игры: wow,", []string{"игры", ":", "wow", ","})

	wT := func(src string, wait string) {
		if WordsOnly(strings.NewReader(src)) != wait {
			t.Fatal("WordsOnly failed for: " + src)
		}
	}

	wT("", "")
	wT(" ", "")
	wT("12", "12")
	wT("12, 13", "12 13")
	wT("12: 13", "12 13")
	wT("Hello, World!", "hello world")
}
