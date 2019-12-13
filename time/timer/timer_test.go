package timer

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	tm := New()

	if tm.Delta() != 0 {
		t.Fatal("not work")
	}

	<-time.After(time.Second)

	if tm.Delta() != 1 {
		t.Fatal("delta not work")
	}
}
