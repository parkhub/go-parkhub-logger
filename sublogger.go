package log

import (
	"fmt"
)

// MARK: Types

// sublogger allows for logging with additional tags
type sublogger struct {
	Logger
	subTags    []string
	skipOffset int
}

// MARK: Public Functions

// Sublogger returns a new sublogger with the provided tags
func Sublogger(tags ...string) Logger {
	return &sublogger{
		Logger:     LoggerSingleton,
		subTags:    tags,
		skipOffset: 0,
	}
}

// Sublogger returns a new sublogger with the provided tags
func (sl *sublogger) Sublogger(tags ...string) Logger {
	return &sublogger{
		Logger:     sl,
		subTags:    tags,
		skipOffset: 1,
	}
}

// MARK: base log methods

// Logln prints the output followed by a newline
func (sl *sublogger) Logln(level Level, message string) {
	sl.printMessage(message, level, nil)
}

// Logf prints the formatted output
func (sl *sublogger) Logf(level Level, format string, a ...interface{}) {
	sl.printMessage(
		fmt.Sprintf(format, a...),
		level,
		nil,
	)
}

// Logd prints output string and data
func (sl *sublogger) Logd(level Level, message string, d interface{}) {
	sl.printMessage(message, level, d)
}

// MARK: Trace

// Traceln prints the output followed by a newline
func (sl *sublogger) Traceln(message string) {
	sl.Logln(LogLevelTrace, message)
}

// Tracef prints the formatted output
func (sl *sublogger) Tracef(format string, a ...interface{}) {
	sl.Logf(LogLevelTrace, format, a...)
}

// Traced prints the output string and data
func (sl *sublogger) Traced(message string, d interface{}) {
	sl.Logd(LogLevelTrace, message, d)
}

// MARK: Debug

// Debugln prints the output followed by a newline
func (sl *sublogger) Debugln(message string) {
	sl.Logln(LogLevelDebug, message)
}

// Debugf prints the formatted output
func (sl *sublogger) Debugf(format string, a ...interface{}) {
	sl.Logf(LogLevelDebug, format, a...)
}

// Debugd prints the output string and data
func (sl *sublogger) Debugd(message string, d interface{}) {
	sl.Logd(LogLevelDebug, message, d)
}

// MARK: Info

// Infoln prints the output followed by a newline
func (sl *sublogger) Infoln(message string) {
	sl.Logln(LogLevelInfo, message)
}

// Infof prints the formatted output
func (sl *sublogger) Infof(format string, a ...interface{}) {
	sl.Logf(LogLevelInfo, format, a...)
}

// Infod prints the output string and data
func (sl *sublogger) Infod(message string, d interface{}) {
	sl.Logd(LogLevelInfo, message, d)
}

// MARK: Warn

// Warnln prints the output followed by a newline
func (sl *sublogger) Warnln(message string) {
	sl.Logln(LogLevelWarn, message)
}

// Warnf prints the formatted output
func (sl *sublogger) Warnf(format string, a ...interface{}) {
	sl.Logf(LogLevelWarn, format, a...)
}

// Warnd prints the output string and data
func (sl *sublogger) Warnd(message string, d interface{}) {
	sl.Logd(LogLevelWarn, message, d)
}

// MARK: Error

// Errorln prints the output followed by a newline
func (sl *sublogger) Errorln(message string) {
	sl.Logln(LogLevelError, message)
}

// Errorf prints the formatted output
func (sl *sublogger) Errorf(format string, a ...interface{}) {
	sl.Logf(LogLevelError, format, a...)
}

// Errord prints the output string and data
func (sl *sublogger) Errord(message string, d interface{}) {
	sl.Logd(LogLevelError, message, d)
}

// MARK: Fatal

// Fatalln prints the output followed by a newline
func (sl *sublogger) Fatalln(message string) {
	sl.Logln(LogLevelFatal, message)
	sl.exit()
}

// Fatalf prints the formatted output
func (sl *sublogger) Fatalf(format string, a ...interface{}) {
	sl.Logf(LogLevelFatal, format, a...)
	sl.exit()
}

// Fatald prints the output string and data
func (sl *sublogger) Fatald(message string, d interface{}) {
	sl.Logd(LogLevelFatal, message, d)
	sl.exit()
}

// MARK: Panic

// Panicln prints the output followed by a newline
func (sl *sublogger) Panicln(message string) {
	sl.Logln(LogLevelPanic, message)
	sl.exit()
}

// Panicf prints the formatted output
func (sl *sublogger) Panicf(format string, a ...interface{}) {
	sl.Logf(LogLevelPanic, format, a...)
	sl.exit()
}

// Panicd prints the output string and data
func (sl *sublogger) Panicd(message string, d interface{}) {
	sl.Logd(LogLevelPanic, message, d)
	sl.exit()
}

// MARK: Private Methods

// newLogMessage creates a new *logMessage
func (sl *sublogger) newLogMessage(output string, level Level, skipOffset int, d interface{}) *logMessage {
	m := sl.Logger.newLogMessage(output, level, sl.skipOffset+skipOffset, d)
	m.Tags = append(m.Tags, sl.subTags...)
	return m
}

// printMessage prints a logMessage to stdout
func (sl *sublogger) printMessage(output string, level Level, d interface{}) {
	if sl.level() > level {
		return
	}

	fmt.Println(sl.newLogMessage(output, level, 0, d).String())
}
