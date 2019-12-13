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

	// speed 3 words per second (180 word per minute)
	val := total / 3
	if val <= 0 {
		val = 1
	}
	return val
}
