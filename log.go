// Package log provides a singular interface to create logs as well as filtering
// them out based on level. It also provides two types of formatting; json or
// pretty. The logger doesn't ship any logs.
package log

import (
	"fmt"
	"os"
)

// LoggerSingleton is the main logging instance.
var LoggerSingleton *logger

// MARK: Setup functions

// SetupLocalLogger is a convenience function for calling SetupLogger with
// pretty formatted logs, colorized output and no tags.
func SetupLocalLogger(level Level) {
	SetupLogger(level, LogFormatPretty, true, true, nil)
}

// SetupCloudLogger is a convenience function for calling SetupLogger with
// JSON formatted logs, non-colorized output and the given tags.
func SetupCloudLogger(level Level, tags []string) {
	SetupLogger(level, LogFormatJSON, false, true, tags)
}

// SetupLogger creates a new logger.
func SetupLogger(level Level, format Format, colorizeOutput bool, logCaller bool, tags []string) {
	if LoggerSingleton != nil {
		// If the logger has already been created, then update its properties
		LoggerSingleton.rawLevel = level
		LoggerSingleton.format = format
		LoggerSingleton.colorizeOutput = colorizeOutput
		LoggerSingleton.tags = tags
		return
	}

	// Setup logger with options.
	LoggerSingleton = &logger{
		rawLevel:       level,
		format:         format,
		colorizeOutput: colorizeOutput,
		logCaller:      logCaller,
		tags:           tags,
	}
}

// MARK: Standard output

// Stdln prints the output followed by a newline.
func Stdln(output string) {
	fmt.Println(output)
}

// Stdf prints the formatted output.
func Stdf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// Stdd prints output string and data.
func Stdd(output string, d interface{}) {
	fmt.Printf("%s %+v", output, d)
}

// MARK: Generic Log

// Logln prints the output followed by a newline
func Logln(level Level, output string) {
	LoggerSingleton.logln(level, output)
}

// Logf prints the formatted output
func Logf(level Level, format string, a ...interface{}) {
	LoggerSingleton.logf(level, format, a...)
}

// Logd prints output string and data
func Logd(level Level, output string, d interface{}) {
	LoggerSingleton.logd(level, output, d)
}

// MARK: Trace

// Traceln prints the output followed by a newline
func Traceln(output string) {
	Logln(LogLevelTrace, output)
}

// Tracef prints the formatted output
func Tracef(format string, a ...interface{}) {
	Logf(LogLevelTrace, format, a...)
}

// Traced prints the output string and data
func Traced(output string, d interface{}) {
	Logd(LogLevelTrace, output, d)
}

// MARK: Debug

// Debugln prints the output followed by a newline.
func Debugln(output string) {
	Logln(LogLevelDebug, output)
}

// Debugf prints the formatted output.
func Debugf(format string, a ...interface{}) {
	Logf(LogLevelDebug, format, a...)
}

// Debugd prints output string and data.
func Debugd(output string, d interface{}) {
	Logd(LogLevelDebug, output, d)
}

// MARK: Info

// Infoln prints the output followed by a newline.
func Infoln(output string) {
	Infod(output, nil)
}

// Infof prints the formatted output.
func Infof(format string, a ...interface{}) {
	Infod(fmt.Sprintf(format, a...), nil)
}

// Infod prints output string and data.
func Infod(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelInfo, d)
}

// MARK: Warn

// Warnln prints the output followed by a newline.
func Warnln(output string) {
	Warnd(output, nil)
}

// Warnf prints the formatted output.
func Warnf(format string, a ...interface{}) {
	Warnd(fmt.Sprintf(format, a...), nil)
}

// Warnd prints output string and data.
func Warnd(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelWarn, d)
}

// MARK: Error

// Errorln prints the output followed by a newline.
func Errorln(output string) {
	Errord(output, nil)
}

// Errorf prints the formatted output.
func Errorf(format string, a ...interface{}) {
	Errord(fmt.Sprintf(format, a...), nil)
}

// Errord prints output string and data.
func Errord(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelError, d)
}

// MARK: Fatal

// Fatalln prints the output followed by a newline and calls os.Exit(1).
func Fatalln(output string) {
	Fatald(output, nil)
}

// Fatalf prints the formatted output.
func Fatalf(format string, a ...interface{}) {
	Fatald(fmt.Sprintf(format, a...), nil)
}

// Fatald prints output string and data.
func Fatald(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelFatal, d)
	os.Exit(1)
}

// MARK: Private Functions

func joinToString(a ...interface{}) string {
	l := len(a)
	if l == 0 {
		return ""
	}
	format := "%v"
	if l > 1 {
		for i := 1; i < l; i++ {
			format += " %v"
		}
	}
	return fmt.Sprintf(format, a...)
}
