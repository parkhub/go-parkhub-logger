package log

import "github.com/ttacon/chalk"

// Level defined the type for a log level.
type Level int

const (
	// LogLevelDebug debug log level.
	LogLevelDebug Level = 0

	// LogLevelInfo info log level.
	LogLevelInfo Level = 1

	// LogLevelWarn warn log level.
	LogLevelWarn Level = 2

	// LogLevelError error log level.
	LogLevelError Level = 3

	// LogLevelFatal fatal log level.
	LogLevelFatal Level = 4
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
