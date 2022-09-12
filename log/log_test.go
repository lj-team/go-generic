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

// go test ./log/ -race
func TestMarshalJSON(t *testing.T) {
	defLog, _ = Open("global=1")

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			InfoJSON().Params(P{
				"RequestID": "123",
				"Entity":    "123",
			}).Message("")
		}()
		go func() {
			defer wg.Done()
			WarnJSON().Message("test message")
		}()
		go func() {
			defer wg.Done()
			ErrorJSON().Params(P{
				"Entity":   789,
				"External": 789,
			}).Send()
		}()
	}
	wg.Wait()
}

func TestLogerStructed(t *testing.T) {
	defLog, _ = Open("global=1")

	InfoParams("test")
	InfoParams("test", 123, 3.5)
	InfoParams("", 123, 3.5, "test2")

	InfoPairs("test")
	InfoPairs("test", 123, 3.5)
	InfoPairs("test", 123, 3.5, "testKey", "testValue")
	InfoPairs("", 123, 3.5, "testKey", "testValue", "wrong")

	InfoParams("test", map[string]any{"foo": "bar", "bar": "foo"})
	// InfoParams("test", 123, 3.5)
	// InfoParams("", 123, 3.5, "test2")
}
