package timer

import (
	"time"
)

type Timer int64

func New() Timer {
	return Timer(time.Now().UnixNano())
}

func (t Timer) DeltaNano() int64 {
	return time.Now().UnixNano() - int64(t)
}

func (t Timer) Delta() int64 {
	return time.Now().Unix() - int64(t)/1000000000
}

func (t Timer) DeltaFloat() float64 {
	return float64(time.Now().UnixNano()-int64(t)) / 1000000000
}
