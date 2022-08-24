package log

import (
	"time"

	"github.com/lj-team/go-generic/time/strftime"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type P map[string]any

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

	params["timestamp"] = strftime.Format("%Y-%m-%d %H:%M:%S", time.Now())
	params["log_level"] = logLevel

	data, _ := json.MarshalToString(&params)

	l.input <- data
}
