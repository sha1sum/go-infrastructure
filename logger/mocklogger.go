package logger
// NewLogMock returns a new mock logger
func NewLogMock(settings Settings) (*MockLog, error) {
	mockLog := &MockLog{}

	var level string

	switch v := settings.Output.(type) {
	case LogglySettings:
		if level != "" {
			level = v.Level
		}
	case Stderr:
		if level != "" {
			level = v.Level
		}
	case Stdout:
		if level != "" {
			level = v.Level
		}
	case Disk:
		if level != "" {
			level = v.Level
		}
	default:
		return mockLog, ErrLogInvalidType
	}

	if level == "debug" {
		mockLog.debug = true
		mockLog.info = true
	}
	if level == "info" {
		mockLog.debug = false
		mockLog.info = true
	}

	mockLog.text = new(Logger)

	return mockLog, nil
}

// MockLog mock object
type MockLog struct {
	text *Logger

	// behavior flags
	isLoggly bool
	debug    bool
	info     bool
}

// Flush inside mock logger
func (l *MockLog) Flush() {}

// Trace inside mock logger
func (l *MockLog) Trace(title string, args ...interface{}) {}

// Debug inside mock logger
func (l *MockLog) Debug(args ...interface{}) {}

// Debugf inside mock logger
func (l *MockLog) Debugf(format string, args ...interface{}) {}

// Info inside mock logger
func (l *MockLog) Info(args ...interface{}) {}

// Infof inside mock logger
func (l *MockLog) Infof(format string, args ...interface{}) {}

// Warn inside mock logger
func (l *MockLog) Warn(args ...interface{}) {}

// Warnf inside mock logger
func (l *MockLog) Warnf(format string, args ...interface{}) {}

// Error inside mock logger
func (l *MockLog) Error(args ...interface{}) {}

// Errorf inside mock logger
func (l *MockLog) Errorf(format string, args ...interface{}) {}

// Fatal inside mock logger
func (l *MockLog) Fatal(args ...interface{}) {}

// Fatalf inside mock logger
func (l *MockLog) Fatalf(format string, args ...interface{}) {}

// Context inside mock logger
func (l *MockLog) Context(fields Fields) ContextualLogger {
	return &MockContext{
		fields: fields,
		logger: l,
	}
}

// MockContext mock
type MockContext struct {
	fields Fields
	logger *MockLog
}

// Debug inside mock logger
func (c *MockContext) Debug(args ...interface{}) {}

// Debugf inside mock logger
func (c *MockContext) Debugf(format string, args ...interface{}) {}

// Info inside mock logger
func (c *MockContext) Info(args ...interface{}) {}

// Infof inside mock logger
func (c *MockContext) Infof(format string, args ...interface{}) {}

// Warn inside mock logger
func (c *MockContext) Warn(args ...interface{}) {}

// Warnf inside mock logger
func (c *MockContext) Warnf(format string, args ...interface{}) {}

// Error inside mock logger
func (c *MockContext) Error(args ...interface{}) {}

// Errorf inside mock logger
func (c *MockContext) Errorf(format string, args ...interface{}) {}

// Fatal inside mock logger
func (c *MockContext) Fatal(args ...interface{}) {}

// Fatalf inside mock logger
func (c *MockContext) Fatalf(format string, args ...interface{}) {}
