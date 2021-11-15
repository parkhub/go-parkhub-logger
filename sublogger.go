package log

import (
	"fmt"
	"os"
)

// MARK: Types

// sublogger allows for logging with additional tags
type sublogger struct {
	*logger
	subTags []string
}

// MARK: Public Functions

// Sublogger returns a new sublogger with the provided tags
func Sublogger(tags ...string) Logger {
	return sublogger{
		logger:  LoggerSingleton,
		subTags: tags,
	}
}

// MARK: Trace

// Traceln prints the output followed by a newline
func (sl sublogger) Traceln(message string) {
	sl.logln(LogLevelTrace, message)
}

// Tracef prints the formatted output
func (sl sublogger) Tracef(format string, a ...interface{}) {
	sl.logf(LogLevelTrace, format, a...)
}

// Traced prints the output string and data
func (sl sublogger) Traced(message string, d interface{}) {
	sl.logd(LogLevelTrace, message, d)
}

// MARK: Debug

// Debugln prints the output followed by a newline
func (sl sublogger) Debugln(message string) {
	sl.logln(LogLevelDebug, message)
}

// Debugf prints the formatted output
func (sl sublogger) Debugf(format string, a ...interface{}) {
	sl.logf(LogLevelDebug, format, a...)
}

// Debugd prints the output string and data
func (sl sublogger) Debugd(message string, d interface{}) {
	sl.logd(LogLevelDebug, message, d)
}

// MARK: Info

// Infoln prints the output followed by a newline
func (sl sublogger) Infoln(message string) {
	sl.logln(LogLevelInfo, message)
}

// Infof prints the formatted output
func (sl sublogger) Infof(format string, a ...interface{}) {
	sl.logf(LogLevelInfo, format, a...)
}

// Infod prints the output string and data
func (sl sublogger) Infod(message string, d interface{}) {
	sl.logd(LogLevelInfo, message, d)
}

// MARK: Warn

// Warnln prints the output followed by a newline
func (sl sublogger) Warnln(message string) {
	sl.logln(LogLevelWarn, message)
}

// Warnf prints the formatted output
func (sl sublogger) Warnf(format string, a ...interface{}) {
	sl.logf(LogLevelWarn, format, a...)
}

// Warnd prints the output string and data
func (sl sublogger) Warnd(message string, d interface{}) {
	sl.logd(LogLevelWarn, message, d)
}

// MARK: Error

// Errorln prints the output followed by a newline
func (sl sublogger) Errorln(message string) {
	sl.logln(LogLevelError, message)
}

// Errorf prints the formatted output
func (sl sublogger) Errorf(format string, a ...interface{}) {
	sl.logf(LogLevelError, format, a...)
}

// Errord prints the output string and data
func (sl sublogger) Errord(message string, d interface{}) {
	sl.logd(LogLevelError, message, d)
}

// MARK: Fatal

// Fatalln prints the output followed by a newline
func (sl sublogger) Fatalln(message string) {
	sl.logln(LogLevelFatal, message)
	os.Exit(1)
}

// Fatalf prints the formatted output
func (sl sublogger) Fatalf(format string, a ...interface{}) {
	sl.logf(LogLevelFatal, format, a...)
	os.Exit(1)
}

// Fatald prints the output string and data
func (sl sublogger) Fatald(message string, d interface{}) {
	sl.logd(LogLevelFatal, message, d)
	os.Exit(1)
}

// MARK: Private Methods

// newLogMessage creates a new *logMessage
func (sl sublogger) newLogMessage(output string, level Level, d interface{}) *logMessage {
	m := sl.logger.newLogMessage(output, level, d)
	m.Tags = append(m.Tags, sl.subTags...)
	return m
}

// printMessage prints a logMessage to stdout
func (sl sublogger) printMessage(output string, level Level, d interface{}) {
	if sl.level > level {
		return
	}

	fmt.Printf(sl.newLogMessage(output, level, d).String()+"\n")
}

// MARK: base log methods

// logln prints the output followed by a newline
func (sl sublogger) logln(level Level, message string) {
	sl.printMessage(message, level, nil)
}

// logf prints the formatted output
func (sl sublogger) logf(level Level, format string, a ...interface{}) {
	sl.printMessage(
		fmt.Sprintf(format, a...),
		level,
		nil,
	)
}

// logd prints output string and data
func (sl sublogger) logd(level Level, message string, d interface{}) {
	sl.printMessage(message, level, d)
}
