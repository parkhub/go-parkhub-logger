package log

import "github.com/ttacon/chalk"

// Level defined the type for a log level.
type Level int

const (
	// LogLevelTrace trace log level
	LogLevelTrace Level = iota

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
)

// MARK: Methods

func (l Level) color() chalk.Color {
	switch l {
	case LogLevelInfo:
		return chalk.Cyan
	case LogLevelWarn:
		return chalk.Yellow
	case LogLevelError:
		return chalk.Red
	case LogLevelFatal:
		return chalk.Magenta
	default:
		return chalk.Black
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
	default:
		return ""
	}
}
