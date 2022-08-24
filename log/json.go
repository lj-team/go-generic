package log

import (
	"time"

	"github.com/lj-team/go-generic/time/strftime"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type p struct {
	Timestamp string `json:"timestamp"`
	LogLevel  string `json:"log_level"`
	P
}

type P struct {
	RequestID string `json:"request_id,omitempty"`
	Entity    any    `json:"entity,omitempty"`
	External  any    `json:"external,omitempty"`
	Message   string `json:"message"`
}

func InfoJSON(params P) {
	defLog.LoggerJSON("info", params)
}

func WarnJSON(params P) {
	defLog.LoggerJSON("warn", params)
}

func ErrorJSON(params P) {
	defLog.LoggerJSON("error", params)
}

func DebugJSON(params P) {
	defLog.LoggerJSON("debug", params)
}

func FinishJSON(params P) {
	defLog.LoggerJSON("info", params)
	defLog.Close()
}

func (l *Log) LoggerJSON(logLevel string, params P) {

	if l == nil {
		return
	}

	if l.end {
		return
	}

	data, _ := json.MarshalToString(&p{
		Timestamp: strftime.Format("%Y-%m-%d %H:%M:%S", time.Now()),
		LogLevel:  logLevel,
		P:         params,
	})

	l.input <- data
}
