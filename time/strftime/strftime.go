package strftime

import (
	"strings"
	"time"
)

var conversion = map[rune]string{
	'B': "January",
	'b': "Jan",
	'm': "01",
	'A': "Monday",
	'a': "Mon",
	'd': "02",
	'H': "15",
	'I': "03",
	'M': "04",
	'S': "05",
	'Y': "2006",
	'y': "06",
	'p': "PM",
	'Z': "MST",
	'z': "-0700",
	'L': ".000",
}

func Format(format string, t time.Time) string {
	retval := make([]byte, 0, len(format))
	for i, ni := 0, 0; i < len(format); i = ni + 2 {
		ni = strings.IndexByte(format[i:], '%')
		if ni < 0 {
			ni = len(format)
		} else {
			ni += i
		}
		retval = append(retval, []byte(format[i:ni])...)
		if ni+1 < len(format) {
			c := format[ni+1]
			if c == '%' {
				retval = append(retval, '%')
			} else {
				if layoutCmd, ok := conversion[rune(c)]; ok {
					retval = append(retval, []byte(t.Format(layoutCmd))...)
				} else {
					retval = append(retval, '%', c)
				}
			}
		} else {
			if ni < len(format) {
				retval = append(retval, '%')
			}
		}
	}
	return string(retval)
}
