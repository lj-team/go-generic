package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/lj-team/go-generic/time/strftime"
)

func InfoParams(message string, params ...any) {
	send("INFO", message, params, false)
}

func WarnParams(message string, params ...any) {
	send("WARN", message, params, false)
}

func ErrorParams(message string, params ...any) {
	send("ERROR", message, params, false)
}

func DebugParams(message string, params ...any) {
	send("DEBUG", message, params, false)
}

func FinishParams(message string, params ...any) {
	send("INFO", message, params, false)
	defLog.Close()
}

//

func InfoPairs(message string, params ...any) {
	send("INFO", message, params, true)
}

func WarnPairs(message string, params ...any) {
	send("WARN", message, params, true)
}

func ErrorPairs(message string, params ...any) {
	send("ERROR", message, params, true)
}

func DebugPairs(message string, params ...any) {
	send("DEBUG", message, params, true)
}

func FinishPairs(message string, params ...any) {
	send("INFO", message, params, true)
	defLog.Close()
}

//

func send(logLevel, message string, params []any, pairs bool) {
	var paramsString []string

	if len(params) != 0 {
		if pairs && len(params)%2 == 0 {
			for i := 0; i < len(params); i += 2 {
				paramsString = append(paramsString, fmt.Sprintf("%v=%v", params[i], params[i+1]))
			}
		} else {
			for _, param := range params {
				paramsString = append(paramsString, fmt.Sprintf("%v", param))
			}
		}
	}

	if message == "" {
		message = "-"
	}
	if len(paramsString) == 0 {
		paramsString = []string{"-"}
	}

	str := fmt.Sprintf(
		`[%s] [%s] [%s] [%s]`,
		strftime.Format("%Y-%m-%d %H:%M:%S", time.Now()),
		logLevel,
		strings.Join(paramsString, ";"),
		message,
	)

	defLog.loggerStructed(str)
}

func (l *Log) loggerStructed(str string) {

	if l == nil {
		return
	}

	if l.end {
		return
	}

	l.input <- str + "\n"
}
