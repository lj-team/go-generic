package answers

import (
	"strings"
	"testing"
)

func TestAnswers(t *testing.T) {

	data := `awesome_frank:
 - Не может быть
 - О, ДА!
 - Правда?
 - Самому страшно
compliment:
 - Спасибо
 - Очень приятно
 - Я тронут
 - Я рад
disable_cats:
 - Никак
 - Вы правда этого хотите?
do_what_you_want:
 - Как скажете
 - Хозяин-барин
`

	a, err := New(strings.NewReader(data))
	if a == nil || err != nil {
		t.Fatal("Load data filed")
	}

	tF := func(code string) {
		res, _ := a.Get(code)
		for _, v := range a.code2texts[code] {
			if v == res {
				return
			}
		}

		t.Fatal("failed")
	}

	tF("do_what_you_want")
	tF("disable_cats")
}
