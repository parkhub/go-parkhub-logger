package log

import "github.com/ttacon/chalk"

// Level defined the type for a log level.
type Level int

const (
	// logLevelUnset represents an unspecified log level
	logLevelUnset Level = iota

	// LogLevelTrace trace log level
	LogLevelTrace

	// LogLevelDebug debug log level.
	LogLevelDebug

	// LogLevelInfo info log level.
	LogLevelInfo

	// LogLevelWarn warn log level.
	LogLevelWarn

	// LogLevelError error log level.
	LogLevelError

	// LogLevelFatal fatal log level.
	LogLevelFatal

	// LogLevelPanic panic log level
	LogLevelPanic
)

// MARK: Methods

func (l Level) color() chalk.Color {
	switch l {
	default:
		return chalk.ResetColor
	case LogLevelTrace:
		return chalk.Blue
	case LogLevelDebug:
		return chalk.Cyan
	case LogLevelInfo:
		return chalk.Green
	case LogLevelWarn:
		return chalk.Yellow
	case LogLevelError, LogLevelPanic:
		return chalk.Red
	case LogLevelFatal:
		return chalk.Magenta
	}
}

// MARK: String interface methods

func (l Level) String() string {
	switch l {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	case LogLevelPanic:
		return "PANIC"
	default:
		return ""
	}
}
