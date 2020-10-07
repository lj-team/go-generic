package log

var defLog *Log

func SetLevel(level string) {
	defLog.SetLevel(level)
}

func GetLevel() string {
	return defLog.GetLevel()
}

func Trace(str ...interface{}) {
	defLog.Trace(str...)
}

func Tracef(format string, str ...interface{}) {
	defLog.Tracef(format, str...)
}

func Debug(str ...interface{}) {
	defLog.Debug(str...)
}

func Debugf(format string, str ...interface{}) {
	defLog.Debugf(format, str...)
}

func Warn(str ...interface{}) {
	defLog.Warn(str...)
}

func Warnf(format string, str ...interface{}) {
	defLog.Warnf(format, str...)
}

func Info(str ...interface{}) {
	defLog.Info(str...)
}

func Infof(format string, str ...interface{}) {
	defLog.Infof(format, str...)
}

func Error(str ...interface{}) {
	defLog.Error(str...)
}

func Errorf(format string, str ...interface{}) {
	defLog.Errorf(format, str...)
}

func Finish(str ...interface{}) {
	defLog.Finish(str...)
}

func Finishf(format string, str ...interface{}) {
	defLog.Finishf(format, str...)
}

func Fatal(str ...interface{}) {
	defLog.Fatal(str...)
}

func Fatalf(format string, str ...interface{}) {
	defLog.Fatalf(format, str...)
}

func Logger(level string, strs ...interface{}) {
	defLog.Logger(level, strs)
}

func Stack(lvl string) {
	defLog.Stack(lvl)
}

func Close() {
	defLog.Close()
	defLog = nil
}
