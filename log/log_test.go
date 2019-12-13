package log

import (
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

	Close()
}
