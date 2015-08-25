package logger

import (
	"runtime"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// Log provides capabilities to log messages both as color formatted message
// to stdout for developers and as JSON formatted messages sent to Logstash.
type Log struct {
	text *logrus.Logger

	// behavior flags
	debug bool
	info  bool
}

//var _ Logger = &Log{}

// Context wraps the standard Logger methods with additional context.
type Context struct {
	fields logrus.Fields
	logger *Log
}

var _ ContextualLogger = &Context{}

// New creates a Logger and accepts a flag for each log level
// Debug and Info level logging can be disabled.
// Log levels Warn, Error, and Fatal are always logged.
// If debug, all calls will also log caller informaton.
//
// The default format of logs will be JSON. Also supports 'text' which is
// a easier format for people to understand if you are logging to stdout.
func New(debug bool, info bool, format string) *Log {
	text := logrus.New()
	text.Level = logrus.DebugLevel

	switch format {
	case "text":
		text.Formatter = &logrus.TextFormatter{ForceColors: true}
	case "json":
		fallthrough
	default:
		text.Formatter = &logrus.JSONFormatter{}
	}

	return &Log{
		text,
		debug,
		info,
	}
}

// Debug logs a message at the Debug level
func (l *Log) Debug(args ...interface{}) {
	if !l.debug {
		return
	}

	l.text.Debugln(appendCallerInfo(args)...)
}

// Debugf logs a printf formatted message at the Debug level
func (l *Log) Debugf(format string, args ...interface{}) {
	if !l.debug {
		return
	}

	l.text.Debugf(format, appendCallerInfo(args)...)
}

// Info logs a message at the Info level
func (l *Log) Info(args ...interface{}) {
	if !l.info {
		return
	}

	if l.debug {
		l.text.Infoln(appendCallerInfo(args)...)
		return
	}

	l.text.Infoln(args...)
}

// Infof logs a printf formatted message at the Info level
func (l *Log) Infof(format string, args ...interface{}) {
	if l.debug {
		f, a := appendCallerInfof(format, args)
		l.text.Infof(f, a...)

		return
	}

	l.text.Infof(format, args...)
}

// Warn logs a message at the Warn level
func (l *Log) Warn(args ...interface{}) {
	if l.debug {
		l.text.Warnln(appendCallerInfo(args)...)
		return
	}

	l.text.Warnln(args...)
}

// Warnf logs a printf formatted message at the Warn level
func (l *Log) Warnf(format string, args ...interface{}) {
	if l.debug {
		l.text.Warnf(format, appendCallerInfo(args)...)
		return
	}

	l.text.Warnf(format, args...)
}

// Error logs a message at the Error level
func (l *Log) Error(args ...interface{}) {
	if l.debug {
		l.text.Errorln(appendCallerInfo(args)...)
		return
	}

	l.text.Errorln(args...)
}

// Errorf logs a printf formatted message at the Error level
func (l *Log) Errorf(format string, args ...interface{}) {
	if l.debug {
		l.text.Errorf(format, appendCallerInfo(args)...)
		return
	}

	l.text.Errorf(format, args...)
}

// Fatal logs a message at the Fatal level
func (l *Log) Fatal(args ...interface{}) {
	if l.debug {
		l.text.Fatalln(appendCallerInfo(args)...)
		return
	}

	l.text.Fatalln(args...)
}

// Fatalf logs a printf formatted message at the Fatal level
func (l *Log) Fatalf(format string, args ...interface{}) {
	if l.debug {
		l.text.Fatalf(format, appendCallerInfo(args)...)
		return
	}

	l.text.Fatalf(format, args...)
}

// Context creates a Context which can either be immediately used or used
// repeatedly within a specific context. Context should be passed an even
// number of key/value values such as:
// logger.Context("foo", 123).Debug("msg")
func (l *Log) Context(fields Fields) ContextualLogger {
	f := logrus.Fields{}

	for k, v := range fields {
		f[k] = v
	}

	if l.debug {
		file, line := getCaller(2)
		fields["caller"] = shortenCaller(file) + ":" + strconv.Itoa(line)
		//fields["line"] = line
	}

	return &Context{
		fields: f,
		logger: l,
	}
}

// Debug logs a message at the Debug level
func (c *Context) Debug(args ...interface{}) {
	c.logger.text.WithFields(c.fields).Debugln(args...)
}

// Debugf logs a printf formatted message at the Debug level
func (c *Context) Debugf(format string, args ...interface{}) {
	c.logger.text.WithFields(c.fields).Debugf(format, args...)
}

// Info logs a message at the Info level
func (c *Context) Info(args ...interface{}) {
	c.logger.text.WithFields(c.fields).Infoln(args...)
}

// Infof logs a printf formatted message at the Info level
func (c *Context) Infof(format string, args ...interface{}) {
	c.logger.text.WithFields(c.fields).Infof(format, args...)
}

// Warn logs a message at the Warn level
func (c *Context) Warn(args ...interface{}) {
	c.logger.text.WithFields(c.fields).Warnln(args...)
}

// Warnf logs a printf formatted message at the Warn level
func (c *Context) Warnf(format string, args ...interface{}) {
	c.logger.text.WithFields(c.fields).Warnf(format, args...)
}

// Error logs a message at the Error level
func (c *Context) Error(args ...interface{}) {
	c.logger.text.WithFields(c.fields).Errorln(args...)
}

// Errorf logs a printf formatted message at the Error level
func (c *Context) Errorf(format string, args ...interface{}) {
	c.logger.text.WithFields(c.fields).Errorf(format, args...)
}

// Fatal logs a message at the Fatal level
func (c *Context) Fatal(args ...interface{}) {
	c.logger.text.WithFields(c.fields).Fatalln(args...)
}

// Fatalf logs a printf formatted message at the Fatal level
func (c *Context) Fatalf(format string, args ...interface{}) {
	c.logger.text.WithFields(c.fields).Fatalf(format, args...)
}

// getCallerInfo returns file and line information for the code that likly logged
func getCaller(depth int) (file string, line int) {
	var ok bool

	_, file, line, ok = runtime.Caller(depth)

	if !ok {
		file = "???"
		line = 0
	}

	return file, line
}

func appendCallerInfof(format string, args ...interface{}) (string, []interface{}) {
	file, line := getCaller(3)
	//args = append(args,file,line)

	return format + " %s:%d", append(args, shortenCaller(file), line)
}

// appendCallerInfo appends the provided arguments with file and line information from
func appendCallerInfo(args ...interface{}) []interface{} {
	file, line := getCaller(3)

	size := len(args)

	result := make([]interface{}, size+1)

	for i := 0; i < size; i++ {
		result[i] = args[i]
	}

	result[size] = shortenCaller(file) + ":" + strconv.Itoa(line)

	return result

}

func shortenCaller(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i:]
			break
		}
		file = short
	}

	return short
}
