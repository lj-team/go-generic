package text

import (
	"io"
)

func ReadingTime(in io.Reader) int {

	total := 0

	Proc(in, func(w string) {

		if skipWrd[w] {
			return
		}

		total++
	})

	// 210 word per minute
	val := (total * 2) / 7
	if val <= 0 {
		val = 1
	}
	return val
}
