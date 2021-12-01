package log

import (
	"fmt"
	"os"
	"time"
)

// MARK: Types

type Logger interface {
	// Trace
	Traceln(string)
	Tracef(string, ...interface{})
	Traced(string, interface{})

	// Debug
	Debugln(string)
	Debugf(string, ...interface{})
	Debugd(string, interface{})

	// Info
	Infoln(string)
	Infof(string, ...interface{})
	Infod(string, interface{})

	// Warn
	Warnln(string)
	Warnf(string, ...interface{})
	Warnd(string, interface{})

	// Error
	Errorln(string)
	Errorf(string, ...interface{})
	Errord(string, interface{})

	// Fatal
	Fatalln(string)
	Fatalf(string, ...interface{})
	Fatald(string, interface{})

	Sublogger(tags...string) Logger
	newLogMessage(string, Level, interface{}) *logMessage
	level() Level
}

// logger is the basic Logger implementation
type logger struct {
	rawLevel       Level
	format         Format
	tags           []string
	colorizeOutput bool
	logCaller      bool
}

// MARK: Public Methods

// Sublogger returns a new sublogger on the logger with the provided tags
func (l logger) Sublogger(tags ...string) Logger {
	return sublogger{
		Logger:  &l,
		subTags: tags,
	}
}

// MARK: Private Methods

// level returns the Logger's Level
func (l logger) level() Level {
	return l.rawLevel
}

// newLogMessage creates a new logMessage
func (l logger) newLogMessage(output string, level Level, d interface{}) *logMessage {
	return newLogMessage(
		l.format,
		l.colorizeOutput,
		l.logCaller,
		newLogTime(time.Now()),
		level,
		l.tags,
		output,
		d,
	)
}

// printMessage prints the message with the given output, level and data. If
// fatal is true, then os.Exit(1) is called after the log has been printed.
func (l logger) printMessage(output string, level Level, d interface{}) {
	if l.rawLevel > level {
		return
	}

	fmt.Printf(l.newLogMessage(output, level, d).String() + "\n")
}

// MARK: base log methods

// logln prints the output followed by a newline
func (l *logger) logln(level Level, message string) {
	l.printMessage(message, level, nil)
}

// logf prints the formatted output
func (l *logger) logf(level Level, format string, a ...interface{}) {
	l.printMessage(
		fmt.Sprintf(format, a...),
		level,
		nil,
	)
}

// logd prints output string and data
func (l *logger) logd(level Level, message string, d interface{}) {
	l.printMessage(message, level, d)
}

// MARK: Trace

// Traceln prints the output followed by a newline
func (l *logger) Traceln(message string) {
	l.logln(LogLevelTrace, message)
}

// Tracef prints the formatted output
func (l *logger) Tracef(format string, a ...interface{}) {
	l.logf(LogLevelTrace, format, a...)
}

// Traced prints the output string and data
func (l *logger) Traced(message string, d interface{}) {
	l.logd(LogLevelTrace, message, d)
}

// MARK: Debug

// Debugln prints the output followed by a newline
func (l *logger) Debugln(message string) {
	l.logln(LogLevelDebug, message)
}

// Debugf prints the formatted output
func (l *logger) Debugf(format string, a ...interface{}) {
	l.logf(LogLevelDebug, format, a...)
}

// Debugd prints the output string and data
func (l *logger) Debugd(message string, d interface{}) {
	l.logd(LogLevelDebug, message, d)
}

// MARK: Info

// Infoln prints the output followed by a newline
func (l *logger) Infoln(message string) {
	l.logln(LogLevelInfo, message)
}

// Infof prints the formatted output
func (l *logger) Infof(format string, a ...interface{}) {
	l.logf(LogLevelInfo, format, a...)
}

// Infod prints the output string and data
func (l *logger) Infod(message string, d interface{}) {
	l.logd(LogLevelInfo, message, d)
}

// MARK: Warn

// Warnln prints the output followed by a newline
func (l *logger) Warnln(message string) {
	l.logln(LogLevelWarn, message)
}

// Warnf prints the formatted output
func (l *logger) Warnf(format string, a ...interface{}) {
	l.logf(LogLevelWarn, format, a...)
}

// Warnd prints the output string and data
func (l *logger) Warnd(message string, d interface{}) {
	l.logd(LogLevelWarn, message, d)
}

// MARK: Error

// Errorln prints the output followed by a newline
func (l *logger) Errorln(message string) {
	l.logln(LogLevelError, message)
}

// Errorf prints the formatted output
func (l *logger) Errorf(format string, a ...interface{}) {
	l.logf(LogLevelError, format, a...)
}

// Errord prints the output string and data
func (l *logger) Errord(message string, d interface{}) {
	l.logd(LogLevelError, message, d)
}

// MARK: Fatal

// Fatalln prints the output followed by a newline
func (l *logger) Fatalln(message string) {
	l.logln(LogLevelFatal, message)
	os.Exit(1)
}

// Fatalf prints the formatted output
func (l *logger) Fatalf(format string, a ...interface{}) {
	l.logf(LogLevelFatal, format, a...)
	os.Exit(1)
}

// Fatald prints the output string and data
func (l *logger) Fatald(message string, d interface{}) {
	l.logd(LogLevelFatal, message, d)
	os.Exit(1)
}
