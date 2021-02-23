package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var LoggerSingleton *logger

// Level defined the type for a log level.
type Level int

// Format visual format of the log message.
type Format string

const (
	// LogLevelDebug debug log level.
	LogLevelDebug Level = 0

	// LogLevelInfo info log level.
	LogLevelInfo Level = 1

	// LogLevelWarn warn log level.
	LogLevelWarn Level = 2

	// LogLevelError error log level.
	LogLevelError Level = 3

	// LogLevelFatal fatal log level.
	LogLevelFatal Level = 4

	// Pretty is a non json formatted log output.
	Pretty Format = "pretty"

	// JSON is a json formatted log output.
	JSON Format = "json"
)

type logger struct {
	Level     Level
	tags      []string
	debugMode bool
	format    Format
}

type logMessage struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Metadata  interface{} `json:"metadata"`
}

// SetupLogger creates a new logger.
func SetupLogger(level Level, tags []string, debugMode bool, format Format) {
	if LoggerSingleton != nil {
		return
	}

	// Setup logger with options.
	LoggerSingleton = &logger{
		Level:     level,
		tags:      tags,
		debugMode: debugMode,
		format:    format,
	}
}

// Stdln prints the output.
func Stdln(output string) {
	fmt.Println(output)
}

// Stdf prints the formatted output.
func Stdf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// Debugln prints the output.
func Debugln(output string) {
	Debugd(output, nil)
}

// Debugd prints output string and data.
func Debugd(output string, d interface{}) {
	buildMessage(output, "DEBUG", d)
}

// Debugf prints the formatted output.
func Debugf(format string, a ...interface{}) {
	Debugln(fmt.Sprintf(format, a...))
}

// Infoln prints the output.
func Infoln(output string) {
	Infod(output, nil)
}

// Infof prints the formatted output.
func Infof(format string, a ...interface{}) {
	Infoln(fmt.Sprintf(format, a...))
}

// Infod prints output string and data.
func Infod(output string, d interface{}) {
	buildMessage(output, "INFO", d)
}

// Warnln prints the output.
func Warnln(output string) {
	Warnd(output, nil)
}

// Warnf prints the formatted output.
func Warnf(format string, a ...interface{}) {
	Warnln(fmt.Sprintf(format, a...))
}

// Warnd prints output string and data.
func Warnd(output string, d interface{}) {
	buildMessage(output, "WARN", d)
}

// Errorln prints the output.
func Errorln(output string) {
	Errord(output, nil)
}

// Errorf prints the formatted output.
func Errorf(format string, a ...interface{}) {
	Errorln(fmt.Sprintf(format, a...))
}

// Errord prints output string and data.
func Errord(output string, d interface{}) {
	buildMessage(output, "ERROR", d)
}

// Fatalln prints the output.
func Fatalln(output string) {
	Fatald(output, nil)
}

// Fatalf prints the formatted output.
func Fatalf(format string, a ...interface{}) {
	Fatalln(fmt.Sprintf(format, a...))
}

// Fatald prints output string and data.
func Fatald(output string, d interface{}) {
	buildMessage(output, "FATAL", d)

}

// MARK: Private

func buildMessage(output string, messageType string, d interface{}) {
	if LoggerSingleton.Level > LogLevelDebug {
		return
	}

	var message *logMessage
	var prettyMessage string

	time := time.Now()
	formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
		time.Year(), time.Month(), time.Day(),
		time.Hour(), time.Minute(), time.Second())

	if LoggerSingleton.format == JSON {
		if d == nil {
			// Format message.
			message = newMessage(formattedTime, messageType, output)
		} else {
			// Format message.
			message = newMessage(formattedTime, messageType, output, d)
		}

		fmt.Println(message.JSONString())
	} else {
		if d == nil {
			// Format message.
			prettyMessage = fmt.Sprintf("%v [%s] %s", formattedTime, messageType, output)
		} else {
			// Format message.
			prettyMessage = fmt.Sprintf("%v [%s] %s %+v", formattedTime, messageType, output, d)
		}

		fmt.Println(prettyMessage)
	}
}

func tagList() string {
	return strings.Join(LoggerSingleton.tags, ",")
}

func newMessage(timestamp string, level string, message string, data ...interface{}) *logMessage {
	formatedMessage := &logMessage{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
		Metadata:  data,
	}

	return formatedMessage
}

// MARK: String methods

func (lm logMessage) JSONString() string {
	s, _ := json.Marshal(lm)
	return string(s)
}

// String string method override.
func (lm logMessage) String() string {
	return lm.JSONString()
}
