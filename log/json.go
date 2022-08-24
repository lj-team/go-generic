package log

import (
	"time"

	"github.com/lj-team/go-generic/time/strftime"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LogJSON struct {
	logLevel string
	params   P
}

type P map[string]any

func newLog(logLevel string) *LogJSON {
	return &LogJSON{
		logLevel: logLevel,
		params:   P{},
	}
}

func InfoJSON() *LogJSON {
	return newLog("info")
}

func WarnJSON() *LogJSON {
	return newLog("warn")
}

func ErrorJSON() *LogJSON {
	return newLog("error")
}

func DebugJSON() *LogJSON {
	return newLog("debug")
}

func (lj *LogJSON) Params(params P) *LogJSON {
	lj.params = params
	return lj
}

func (lj *LogJSON) Message(msg string) {
	if msg != "" {
		lj.params["message"] = msg
	}
	lj.send()
}

func (lj *LogJSON) Finish(msg string) {
	if msg != "" {
		lj.params["message"] = msg
	}
	lj.send()
	defLog.Close()
}

func (lj *LogJSON) Send() {
	lj.send()
}

func (lj *LogJSON) send() {
	lj.params["timestamp"] = strftime.Format("%Y-%m-%d %H:%M:%S", time.Now())
	lj.params["log_level"] = lj.logLevel

	defLog.loggerJSON(lj.params)
}

func (l *Log) loggerJSON(params P) {

	if l == nil {
		return
	}

	if l.end {
		return
	}

	data, _ := json.MarshalToString(&params)

	l.input <- data + "\n"
}
