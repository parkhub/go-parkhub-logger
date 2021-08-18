package log

import (
	"fmt"
	"time"
)

type logger struct {
	level          Level
	format         Format
	tags           []string
	colorizeOutput bool
	logCaller      bool
}

// MARK: Private methods

// printMessage prints the message with the given output, level and data. If
// fatal is true, then os.Exit(1) is called after the log has been printed.
func (l logger) printMessage(output string, level Level, d interface{}) {
	if l.level > level {
		return
	}

	fmt.Printf(
		newLogMessage(
			l.format,
			l.colorizeOutput,
			l.logCaller,
			newLogTime(time.Now()),
			level,
			l.tags,
			output,
			d,
		).String()+"\n",
	)
}
