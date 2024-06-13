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
	// According to loggly documentation:
	// * The only timestamp format accepted is ISO 8601 (e.g., 2013-10-11T22:14:15.003Z).
	// * Loggly supports microseconds/seconds fraction up to 6 digits, per the spec in RFC5424.

	TimeFormatDefault TimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
	TimeFormatLoggly  TimeFormat = "2006-01-02T15:04:05.000000Z07:00"
)
