package log

// Format visual format of the log message.
type Format string

const (
	// LogFormatPretty is a non json formatted log output.
	LogFormatPretty Format = "pretty"

	// LogFormatJSON is a json formatted log output.
	LogFormatJSON Format = "json"
)

type TimeFormat string

const (
	TimeFormatLoggly TimeFormat = "loggly"
)
