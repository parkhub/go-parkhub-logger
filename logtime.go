package log

import (
	"fmt"
	"math"
	"time"
)

// logTime converts a time to a string suitable for logging.
type logTime struct {
	// The log's time.
	time time.Time

	// The log's centiseconds. Suitable for logging with the time format
	// 00-00-00T00:00:00-XX:00 where XX represents the 1/100th of a second this
	// property stores.
	centiseconds int

	// The log's centiseconds remainder. Suitable for logging with the time format
	// 00-00-00T00:00:00-00:XX where XX represents the 1/100th of a centisecond
	// this property stores.
	centisecondsRemainder int
}

// MARK: Initializers

// newLogTime creates and returns a new log time suitible for logging
// statements.
func newLogTime(time time.Time) logTime {
	nanoseconds := float64(time.Nanosecond())
	centiseconds := math.Floor(float64(nanoseconds) / 1.0e7)
	nanoRemainder := nanoseconds - centiseconds*1.0e7
	centisecondRemainder := math.Floor(float64(nanoRemainder) / 1.0e5)

	return logTime{
		time:                  time,
		centiseconds:          int(centiseconds),
		centisecondsRemainder: int(centisecondRemainder),
	}
}

// MARK: String interface methods

func (t logTime) String() string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-%02d:%02d",
		t.time.Year(), t.time.Month(), t.time.Day(),
		t.time.Hour(), t.time.Minute(), t.time.Second(),
		t.centiseconds, t.centisecondsRemainder)
}
