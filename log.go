// Package log provides a singular interface to create logs as well as filtering
// them out based on level. It also provides two types of formatting; json or
// pretty. The logger doesn't ship any logs.
package log

import (
	"fmt"
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
		LoggerSingleton.level = level
		LoggerSingleton.format = format
		LoggerSingleton.colorizeOutput = colorizeOutput
		LoggerSingleton.tags = tags
	}

	// Setup logger with options.
	LoggerSingleton = &logger{
		level:          level,
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
	fmt.Printf(fmt.Sprintf("%s %+v", output, d))
}

// Stddln prints output string and data followed by a newline.
func Stddln(output string, d interface{}) {
	Stdd(output+"\n", d)
}

// MARK: Debug

// Debugln prints the output followed by a newline.
func Debugln(output string) {
	Debugd(output+"\n", nil)
}

// Debugf prints the formatted output.
func Debugf(format string, a ...interface{}) {
	Debugd(fmt.Sprintf(format, a...), nil)
}

// Debugd prints output string and data.
func Debugd(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelDebug, false, d)
}

// Debugdln prints output string and data followed by a newline.
func Debugdln(output string, d interface{}) {
	Debugd(output+"\n", d)
}

// MARK: Info

// Infoln prints the output followed by a newline.
func Infoln(output string) {
	Infod(output+"\n", nil)
}

// Infof prints the formatted output.
func Infof(format string, a ...interface{}) {
	Infod(fmt.Sprintf(format, a...), nil)
}

// Infod prints output string and data.
func Infod(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelInfo, false, d)
}

// Infodln prints output string and data followed by a newline.
func Infodln(output string, d interface{}) {
	Infod(output+"\n", d)
}

// MARK: Warn

// Warnln prints the output followed by a newline.
func Warnln(output string) {
	Warnd(output+"\n", nil)
}

// Warnf prints the formatted output.
func Warnf(format string, a ...interface{}) {
	Warnd(fmt.Sprintf(format, a...), nil)
}

// Warnd prints output string and data.
func Warnd(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelWarn, false, d)
}

// Warndln prints output string and data followed by a newline.
func Warndln(output string, d interface{}) {
	Warnd(output+"\n", d)
}

// MARK: Error

// Errorln prints the output followed by a newline.
func Errorln(output string) {
	Errord(output+"\n", nil)
}

// Errorf prints the formatted output.
func Errorf(format string, a ...interface{}) {
	Errord(fmt.Sprintf(format, a...), nil)
}

// Errord prints output string and data.
func Errord(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelError, false, d)
}

// Errordln prints output string and data followed by a newline.
func Errordln(output string, d interface{}) {
	Errord(output+"\n", d)
}

// MARK: Fatal

// Fatalln prints the output followed by a newline and calls os.Exit(1).
func Fatalln(output string) {
	Fatald(output+"\n", nil)
}

// Fatalf prints the formatted output.
func Fatalf(format string, a ...interface{}) {
	Fatald(fmt.Sprintf(format, a...), nil)
}

// Fatald prints output string and data.
func Fatald(output string, d interface{}) {
	LoggerSingleton.printMessage(output, LogLevelFatal, true, d)
}

// Fataldln prints output string and data.
func Fataldln(output string, d interface{}) {
	Fatald(output+"\n", d)
}
