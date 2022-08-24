package log

import (
	"sync"
	"testing"
)

func TestLoggerSetLevel(t *testing.T) {

	defLog, _ = Open("global=1")

	if defLog == nil {
		t.Fatal("defLog=nil")
	}

	if GetLevel() != "info" {
		t.Fatal("invalid default log level")
	}

	SetLevel("error")

	if GetLevel() != "error" {
		t.Fatal("expected log level 'error'")
	}

	Errorf("failed %d", 1)
	Finishf("finish %s", "ok")
	Stack("error")

	defLog.Close()

	defLog = nil

	Error("error1", "error2")
	Warn("warn")
	Info("info")
	Debug("debug")
	Trace("trace", 1, "trace", 2)
	Infof("Hello, %s", "Mike")
	Errorf("Hello, %s", "Mike")
	Warnf("Hello, %s", "Mike")
	Debugf("Hello, %s", "Mike")
	Tracef("Hello, %s", "Mike")
	Stack("error")

	Close()
}

func TestMarshalJSON(t *testing.T) {
	defLog, _ = Open("global=1")

	p1 := P{
		RequestID: "123",
		Entity:    "123",
		External:  123,
		Message:   "Test1",
	}
	p2 := P{
		RequestID: "456",
		Entity:    "456",
		External:  456,
		Message:   "Test2",
	}
	p3 := P{
		RequestID: "789",
		Entity:    789,
		External:  "789",
		Message:   "Test3",
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			InfoJSON(p1)
		}()
		go func() {
			defer wg.Done()
			WarnJSON(p2)
		}()
		go func() {
			defer wg.Done()
			ErrorJSON(p3)
		}()
	}
	wg.Wait()
}
