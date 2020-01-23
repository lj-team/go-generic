package time

import (
	"time"
)

func DaysInMonth(year int, month int) int {

	if month != 12 {
		tm := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
		d := tm.Day()
		return d
	}

	return 31
}
