package time

import (
	"testing"
)

func TestDaysInMonth(t *testing.T) {

	tF := func(y, m, d int) {
		if DaysInMonth(y, m) != d {
			t.Fatalf("DaysInMonth faild for y=%d m=%d", y, m)
		}
	}

	tF(2019, 1, 31)
	tF(2019, 2, 28)
	tF(2019, 3, 31)
	tF(2019, 4, 30)
	tF(2019, 5, 31)
	tF(2019, 6, 30)
	tF(2019, 7, 31)
	tF(2019, 8, 31)
	tF(2019, 9, 30)
	tF(2019, 10, 31)
	tF(2019, 11, 30)
	tF(2019, 12, 31)
	tF(2020, 1, 31)
	tF(2020, 2, 29)
	tF(2020, 3, 31)
	tF(2020, 4, 30)
	tF(2020, 5, 31)
	tF(2020, 6, 30)
	tF(2020, 7, 31)
	tF(2020, 8, 31)
	tF(2020, 9, 30)
	tF(2020, 10, 31)
	tF(2020, 11, 30)
	tF(2020, 12, 31)
}
