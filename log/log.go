package log

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/lj-team/go-generic/text/kv"
	"github.com/lj-team/go-generic/time/strftime"
)

const (
	eof = "9e93c0c16d5bb7447e68e1d15e64215e"
)

var logLevels map[string]int = map[string]int{
	"none":  0,
	"fatal": 1,
	"error": 2,
	"warn":  3,
	"info":  4,
	"debug": 5,
	"trace": 6,
}

type LoggerFunc func(str ...interface{})

type Log struct {
	level     int
	conf      *Config
	input     chan string
	fh        *os.File
	filename  string
	lastCheck int64
	eofC      chan bool
	end       bool
}

func (l *Log) Logger(level string, strs []interface{}) {

	if l == nil {
		return
	}

	if l.end {
		return
	}

	code, ok := logLevels[level]
	if ok && code <= l.level {
		for _, text := range strs {
			l.input <- strftime.Format("%Y-%m-%d %H:%M:%S", time.Now()) + " | " + level + " | " + fmt.Sprint(text) + "\n"
		}
	}
}

func New(c *Config, def bool) (*Log, error) {

	if def && defLog != nil {
		return defLog, nil
	}

	l := &Log{
		conf:      c,
		filename:  strftime.Format(c.Template, time.Now()),
		lastCheck: time.Now().Unix(),
		input:     make(chan string, 1024),
		eofC:      make(chan bool),
	}

	if c.Template != "" {
		var err error

		if l.fh, err = os.OpenFile(l.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err != nil {
			return nil, err
		}

		if c.Save > 0 {
			rm_name := strftime.Format(c.Template, time.Unix(l.lastCheck-int64(c.Save*c.Period), 0))
			os.Remove(rm_name)
		}
	}

	l.SetLevel(c.Level)

	if def && defLog == nil {
		defLog = l
	}

	go l.writer()

	return l, nil
}

func Init(c *Config) {
	New(c, true)
}

func Open(params string) (*Log, error) {

	args, err := kv.New(params)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.StdErr = args.GetBool("stderr", false)
	cfg.StdOut = args.GetBool("stdout", false)

	name := args.GetString("name", "")
	if name != "" {
		template := args.GetString("path", ".") + "/" + name
		period := args.GetString("period", "day")

		switch period {
		case "day":
			template += "-%Y%m%d.log"
			cfg.Period = 86400
		case "hour":
			template += "-%Y%m%d%H.log"
			cfg.Period = 3600
		case "month":
			template += "-%Y%m.log"
			cfg.Period = 86400 * 31
		default:
			template += "-%Y%m%d.log"
			cfg.Period = 86400
		}

		cfg.Template = template
	}

	cfg.Save = args.GetInt("save", 14)
	cfg.Level = args.GetString("level", "info")

	return New(cfg, args.GetBool("global", true))
}

func (l *Log) rotate() {

	if l.lastCheck+60 > time.Now().Unix() {
		return
	}

	l.lastCheck = time.Now().Unix()
	new_name := strftime.Format(l.conf.Template, time.Now())

	if new_name != l.filename {
		l.fh.Close()

		var err error
		l.filename = new_name

		if l.fh, err = os.OpenFile(l.filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err != nil {
			panic(err)
		}

		if l.conf.Save > 0 {
			rm_name := strftime.Format(l.conf.Template, time.Unix(l.lastCheck-int64(l.conf.Save*l.conf.Period), 0))
			os.Remove(rm_name)
		}
	}
}

func (l *Log) writer() {
	for {
		select {

		case str := <-l.input:

			conf := l.conf

			if str == eof {

				if l.fh != nil {
					l.fh.Close()
				}

				close(l.eofC)
				return
			}

			if l.fh != nil {
				l.rotate()
				l.fh.WriteString(str)
				l.fh.Sync()
			}

			if conf.StdOut {
				os.Stdout.WriteString(str)
				os.Stdout.Sync()
			}

			if conf.StdErr {
				os.Stderr.WriteString(str)
				os.Stderr.Sync()
			}

		case <-time.After(time.Minute):

			if l.fh != nil {
				l.rotate()
			}
		}
	}
}

func (l *Log) Close() {
	if l != nil {
		if !l.end {
			l.end = true
			l.input <- eof
			<-l.eofC
			close(l.input)
		}
	}
}

func (l *Log) Fatal(str ...interface{}) {
	l.Logger("fatal", str)
	l.Close()
	os.Exit(1)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func (l *Log) Finish(str ...interface{}) {
	l.Logger("info", str)
	l.Close()
}

func (l *Log) Finishf(format string, args ...interface{}) {
	l.Finish(fmt.Sprintf(format, args...))
}

func (l *Log) Error(str ...interface{}) {
	l.Logger("error", str)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Log) Info(str ...interface{}) {
	l.Logger("info", str)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Log) Debug(str ...interface{}) {
	l.Logger("debug", str)
}

func (l *Log) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Log) Warn(str ...interface{}) {
	l.Logger("warn", str)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Log) Trace(str ...interface{}) {
	l.Logger("trace", str)
}

func (l *Log) Tracef(format string, args ...interface{}) {
	l.Trace(fmt.Sprintf(format, args...))
}

func (l *Log) SetLevel(lvl string) {
	if l == nil {
		return
	}

	if code, ok := logLevels[lvl]; ok {
		l.level = code
	} else {
		l.level = 0
	}
}

func (l *Log) GetLevel() string {
	if l == nil {
		return "none"
	}

	for code, level := range logLevels {
		if level == l.level {
			return code
		}
	}
	return "none"
}

func (l *Log) Stack(lvl string) {
	buf := make([]byte, 1024*100)
	n := runtime.Stack(buf, false)
	l.Logger(lvl, []interface{}{string(buf[:n])})
}
