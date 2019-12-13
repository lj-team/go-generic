package wdog

import (
	"testing"
	"time"
)

func TestWatchDog(t *testing.T) {

	msg := ""

	wd := New(1, func() {
		msg = "wdog"
	})
	defer wd.Close()

	for i := 0; i < 25; i++ {
		wd.Alive()
		<-time.After(time.Nanosecond * 100000000)
		if msg != "" {
			t.Fatal("wdog kill")
		}
	}

	<-time.After(time.Second * 2)
	if msg == "" {
		t.Fatal("wdog not work")
	}
}
