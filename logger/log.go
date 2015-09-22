package logger

import (
	"errors"
	"os"
	"runtime"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/sebest/logrusly"
)

// Settings is container for an interface that is represented by the various
// types below.
type Settings struct {
	Output interface{}
	Debug  bool
	Info   bool
}

// LogglySettings is a type of output using the Loggly service.
type LogglySettings struct {
	Level  string
	Token  string
	Domain string
	Tags   []string
}

// Stderr is a type of output that uses the os.Stderr
type Stderr struct {
	Level  string
	Format string
}

// Stdout is a type of output that uses the os.Stdout
type Stdout struct {
	Level  string
	Format string
}

// Disk is a type of output that uses the logrus output which writes to disk.
type Disk struct {
	Path  string
	Level string
}

// Log provides capabilities to log messages both as color formatted message
// to stdout for developers and as JSON formatted messages sent to Logstash.
type Log struct {
	text *logrus.Logger
	hook *logrusly.LogglyHook

	isLoggly bool
	debug    bool
	info     bool
}

// Context wraps the standard Logger methods with additional context.
type Context struct {
	fields logrus.Fields
	logger *Log
}

var (
	// ErrLogLogglySetup is a error that informs that we're missing a token.
	ErrLogLogglySetup = errors.New("Please make sure your loggly token / domain is setup properly")
	// ErrLogInvalidLevel is an error that is thrown when an invalid string is
	// passed as a level.
	ErrLogInvalidLevel = errors.New("Please make sure you use a valid level: panic, fatal, error, warn, info, debug")
	// ErrLogInvalidPath is an error that is thrown when an invalid path is passed
	// into Disk
	ErrLogInvalidPath = errors.New("Invalid path or permission error.")
	// ErrLogInvalidType is an error that is thrown when settings.Output.(type)
	// doesn't match what we're expecting.
	ErrLogInvalidType = errors.New("Invalid log output type.")
)

// New creates a Logger
// Log levels Warn, Error, and Fatal are always logged.
// If debug, all calls will also log caller informaton.
//
// The default format of logs will be JSON. Also supports 'text' which is
// a easier format for people to understand if you are logging to stdout.
func New(settings Settings) (*Log, error) {
	// Setting up Log.
	log := &Log{}

	// Let's fire up a new logrus
	text := logrus.New()
	text.Level = logrus.DebugLevel

	// Default is text, color me pretty if possible.
	text.Formatter = &logrus.TextFormatter{ForceColors: true}

	var level logrus.Level
	var err error
	var overrideDefaultLevel bool

	switch v := settings.Output.(type) {
	case LogglySettings:
		// Validate that the bare minimum of what is needed is present.
		if v.Token == "" && v.Domain == "" {
			return &Log{}, ErrLogLogglySetup
		}

		if v.Level != "" && v.Level != "debug" {
			level, err = logrus.ParseLevel(v.Level)
			if err != nil {
				return log, ErrLogInvalidLevel
			}
			overrideDefaultLevel = true
		}

		var token, domain string
		// Change the formatter to JSON, required here.
		text.Formatter = &logrus.JSONFormatter{}
		// initialize LogglyHook
		token = v.Token
		domain = v.Domain
		// Create new hook
		hook := logrusly.NewLogglyHook(token, domain, level)
		// Add hook tags
		for _, tag := range v.Tags {
			hook.Tag(tag)
		}
		// Finally, add the hook to the logrus
		text.Hooks.Add(hook)
		// Let's add information to the main log.
		log.hook = hook
		log.isLoggly = true
	case Stderr:
		text.Out = os.Stderr

		if v.Format == "json" {
			text.Formatter = &logrus.JSONFormatter{}
		}

		if v.Level != "" && v.Level != "debug" {
			level, err = logrus.ParseLevel(v.Level)
			if err != nil {
				return log, ErrLogInvalidLevel
			}
			overrideDefaultLevel = true
		}
	case Stdout:
		text.Out = os.Stdout

		if v.Format == "json" {
			text.Formatter = &logrus.JSONFormatter{}
		}

		if v.Level != "" && v.Level != "debug" {
			level, err = logrus.ParseLevel(v.Level)
			if err != nil {
				return log, ErrLogInvalidLevel
			}
			overrideDefaultLevel = true
		}
	case Disk:
		if v.Path == "" {
			return log, ErrLogInvalidPath
		}

		f, err := os.OpenFile(v.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return log, ErrLogInvalidPath
		}

		text.Out = f
		text.Formatter = &logrus.JSONFormatter{}

		if v.Level != "" && v.Level != "debug" {
			level, err = logrus.ParseLevel(v.Level)
			if err != nil {
				return log, err
			}
			overrideDefaultLevel = true
		}
	default:
		// We should assume text unless overriden otherwise.
		return log, ErrLogInvalidType
	}

	// All the error checking for level has been handled up above.
	if overrideDefaultLevel {
		text.Level = level
	}

	// Log debug / info evaluation
	if text.Level == logrus.DebugLevel {
		log.debug = true
		log.info = true
	}

	if text.Level == logrus.InfoLevel {
		log.debug = false
		log.info = true
	}

	// Let's attach text to the log.
	log.text = text

	return log, nil
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

// Flush is a pass-through function that is exposed on logrusly package.
func (l *Log) Flush() {
	if l.isLoggly {
		// If this is loggly, we're passing flush to the client otherwise we lose
		// 5 seconds worth of data (for not panic/fatal messages) and that's not cool.
		l.hook.Flush()
	}
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
