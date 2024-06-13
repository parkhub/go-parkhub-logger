package log

import (
	"fmt"
	"time"
)

// MARK: Types

type Logger interface {
	// Base log
	Logln(Level, string)
	Logf(Level, string, ...interface{})
	Logd(Level, string, interface{})

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

	// Create a logger object with additional tags
	Sublogger(tags ...string) Logger

	// Recover from a panic, log it, and set an error variable
	Recover(label string, err *error)

	// Private methods
	newLogMessage(message string, level Level, skipOffset int, data interface{}) *logMessage
	level() Level
	exit()
}

// logger is the basic Logger implementation
type logger struct {
	rawLevel       Level
	format         Format
	timeFormat     TimeFormat
	tags           []string
	colorizeOutput bool
	logCaller      bool
	exitFunc       func()
}

// MARK: Public Methods

// Sublogger returns a new sublogger on the logger with the provided tags
func (l *logger) Sublogger(tags ...string) Logger {
	return &sublogger{
		Logger:  l,
		subTags: tags,
	}
}

// MARK: Private Methods

// level returns the Logger's Level
func (l *logger) level() Level {
	return l.rawLevel
}

// newLogMessage creates a new logMessage
func (l *logger) newLogMessage(output string, level Level, skipOffset int, d interface{}) *logMessage {
	return newLogMessage(
		l.format,
		l.colorizeOutput,
		l.logCaller,
		5+skipOffset,
		time.Now(),
		l.timeFormat,
		level,
		l.tags,
		output,
		d,
	)
}

// printMessage prints the message with the given output, level and data. If
// fatal is true, then os.Exit(1) is called after the log has been printed.
func (l *logger) printMessage(output string, level Level, d interface{}) {
	if l.rawLevel > level {
		return
	}

	fmt.Println(l.newLogMessage(output, level, 0, d).String())
}

func (l *logger) exit() {
	l.exitFunc()
}

// MARK: base log methods

// Logln prints the output followed by a newline
func (l *logger) Logln(level Level, message string) {
	l.printMessage(message, level, nil)
}

// Logf prints the formatted output
func (l *logger) Logf(level Level, format string, a ...interface{}) {
	l.printMessage(
		fmt.Sprintf(format, a...),
		level,
		nil,
	)
}

// Logd prints output string and data
func (l *logger) Logd(level Level, message string, d interface{}) {
	l.printMessage(message, level, d)
}

// MARK: Trace

// Traceln prints the output followed by a newline
func (l *logger) Traceln(message string) {
	l.Logln(LogLevelTrace, message)
}

// Tracef prints the formatted output
func (l *logger) Tracef(format string, a ...interface{}) {
	l.Logf(LogLevelTrace, format, a...)
}

// Traced prints the output string and data
func (l *logger) Traced(message string, d interface{}) {
	l.Logd(LogLevelTrace, message, d)
}

// MARK: Debug

// Debugln prints the output followed by a newline
func (l *logger) Debugln(message string) {
	l.Logln(LogLevelDebug, message)
}

// Debugf prints the formatted output
func (l *logger) Debugf(format string, a ...interface{}) {
	l.Logf(LogLevelDebug, format, a...)
}

// Debugd prints the output string and data
func (l *logger) Debugd(message string, d interface{}) {
	l.Logd(LogLevelDebug, message, d)
}

// MARK: Info

// Infoln prints the output followed by a newline
func (l *logger) Infoln(message string) {
	l.Logln(LogLevelInfo, message)
}

// Infof prints the formatted output
func (l *logger) Infof(format string, a ...interface{}) {
	l.Logf(LogLevelInfo, format, a...)
}

// Infod prints the output string and data
func (l *logger) Infod(message string, d interface{}) {
	l.Logd(LogLevelInfo, message, d)
}

// MARK: Warn

// Warnln prints the output followed by a newline
func (l *logger) Warnln(message string) {
	l.Logln(LogLevelWarn, message)
}

// Warnf prints the formatted output
func (l *logger) Warnf(format string, a ...interface{}) {
	l.Logf(LogLevelWarn, format, a...)
}

// Warnd prints the output string and data
func (l *logger) Warnd(message string, d interface{}) {
	l.Logd(LogLevelWarn, message, d)
}

// MARK: Error

// Errorln prints the output followed by a newline
func (l *logger) Errorln(message string) {
	l.Logln(LogLevelError, message)
}

// Errorf prints the formatted output
func (l *logger) Errorf(format string, a ...interface{}) {
	l.Logf(LogLevelError, format, a...)
}

// Errord prints the output string and data
func (l *logger) Errord(message string, d interface{}) {
	l.Logd(LogLevelError, message, d)
}

// MARK: Fatal

// Fatalln prints the output followed by a newline
func (l *logger) Fatalln(message string) {
	l.Logln(LogLevelFatal, message)
	l.exitFunc()
}

// Fatalf prints the formatted output
func (l *logger) Fatalf(format string, a ...interface{}) {
	l.Logf(LogLevelFatal, format, a...)
	l.exitFunc()
}

// Fatald prints the output string and data
func (l *logger) Fatald(message string, d interface{}) {
	l.Logd(LogLevelFatal, message, d)
	l.exitFunc()
}

// MARK: Panic

// Panicln prints the output followed by a newline
func (l *logger) Panicln(message string) {
	l.Logln(LogLevelPanic, message)
	l.exitFunc()
}

// Panicf prints the formatted output
func (l *logger) Panicf(format string, a ...interface{}) {
	l.Logf(LogLevelPanic, format, a...)
	l.exitFunc()
}

// Panicd prints the output string and data
func (l *logger) Panicd(message string, d interface{}) {
	l.Logd(LogLevelPanic, message, d)
	l.exitFunc()
}
