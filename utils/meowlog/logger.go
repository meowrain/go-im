package meowlog

type LogLevel uint

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger interface {
	// enable returns true if the given level is enabled for logging.
	enable(level LogLevel) bool
	// SetLogLevel sets the log level for the logger.
	SetLogLevel(level string)
	// GetLogLevel returns the log level for the given level.
	GetLogLevel(level LogLevel) string
	// Debug logs a message at the DebugLevel.
	Debug(format string, args ...any)
	// Trace logs a message at the TraceLevel.
	Trace(format string, args ...any)
	// Info logs a message at the InfoLevel.
	Info(format string, args ...any)
	// Warn logs a message at the WarnLevel.
	Warn(format string, args ...any)
	// Error logs a message at the ErrorLevel.
	Error(format string, args ...any)
	// Fatal logs a message at the FatalLevel and exits the program.
	Fatal(format string, args ...any)
	// log logs a message at the given level.
	log(level LogLevel, format string, args ...any)
}

func NewLogger(logType string, level string, dir ...string) Logger {
	switch logType {
	case "console":
		return newMeowConsoleLog(level)
	case "file":
		return newMeowFileLog(level, dir[0])
	default:
		return newMeowConsoleLog(level)
	}
}
