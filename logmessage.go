package log

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"unicode"

	"github.com/ttacon/chalk"
)

type logMessage struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Tags      []string    `json:"tags"`
	Message   string      `json:"message"`
	Metadata  interface{} `json:"metadata,omitempty"`
	File      string      `json:"file,omitempty"`

	format       Format
	rawLevel     Level
	colorize     bool
	logCaller    bool
	trimmedLeft  string
	trimmedRight string
}

// newLogMessage creates a new logMessage.
func newLogMessage(
	format Format,
	colorize bool,
	logCaller bool,
	callerSkip int,
	time logTime,
	level Level,
	tags []string,
	message string,
	data interface{},
) *logMessage {
	trimmedLeft := leadingWhitespace(message)
	trimmedRight := trailingWhitespace(message)
	modifiedMessage := strings.TrimSpace(message)

	caller := ""
	if logCaller {
		var file string
		var line int
		var ok bool
		_, file, line, ok = runtime.Caller(callerSkip + 1)
		if ok {
			fileComponents := strings.Split(file, "/")
			if len(fileComponents) > 1 {
				file = fileComponents[len(fileComponents)-1]
			}
		} else {
			file = "???"
			line = 0
		}
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	// Make sure that error interface types in logMessage.Metadata are not
	// marshalled as empty JSON objects
	metadata := data
	if e, ok := data.(error); ok {
		if _, ok := e.(json.Marshaler); !ok {
			metadata = e.Error()
		}
	}
	// If logMessage.Metadata is a slice, convert make the same conversion for any
	// errors in the slice
	if s, ok := data.([]interface{}); ok {
		for i, v := range s {
			if e, ok := v.(error); ok {
				if _, ok := e.(json.Marshaler); !ok {
					s[i] = e.Error()
				}
			}
		}
	}

	formattedMessage := &logMessage{
		Timestamp:    time.String(),
		Level:        level.String(),
		Tags:         tags,
		Message:      modifiedMessage,
		Metadata:     metadata,
		File:         caller,
		format:       format,
		rawLevel:     level,
		colorize:     colorize,
		logCaller:    logCaller,
		trimmedLeft:  trimmedLeft,
		trimmedRight: trimmedRight,
	}

	return formattedMessage
}

// MARK: Methods

func (m logMessage) restoreWhitespace(output string) string {
	return m.trimmedLeft + output + m.trimmedRight
}

func (m logMessage) colorizeIfNeeded(output string) string {
	if !m.colorize {
		return output
	}

	lines := strings.Split(output, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%s%s%s", m.rawLevel.color(), line, chalk.ResetColor)
	}

	return strings.Join(lines, "\n")
}

func (m logMessage) jsonString() string {
	s, _ := json.Marshal(m)
	return m.colorizeIfNeeded(string(s))
}

// MARK: String interface methods

func (m logMessage) String() string {
	if m.format == LogFormatJSON {
		return m.jsonString()
	}

	var timeAndLevel, logCaller, tags, data string

	timeAndLevel = fmt.Sprintf("%s [%s] ", m.Timestamp, m.Level)

	if m.logCaller {
		logCaller = fmt.Sprintf("[%s] ", m.File)
	}

	if len(m.Tags) > 0 {
		tags = fmt.Sprintf("(%s) ", strings.Join(m.Tags, ","))
	}

	if m.Metadata != nil {
		switch m.Metadata.(type) {
		case string:
			data = m.Metadata.(string)
		case fmt.Stringer:
			data = m.Metadata.(fmt.Stringer).String()
		case error:
			data = m.Metadata.(error).Error()
		default:
			data = fmt.Sprintf("%#v", m.Metadata)
		}
	}

	prettyMessage := timeAndLevel + logCaller + tags +
		m.trimmedLeft + m.Message + " " + data + m.trimmedRight

	return m.colorizeIfNeeded(prettyMessage)
}

// MARK: Helper Functions

// return the leading whitespace of the input string
func leadingWhitespace(s string) string {
	var b strings.Builder
	for i, _ := range s {
		if !unicode.IsSpace(rune(s[i])) {
			return b.String()
		}
		b.Write([]byte{s[i]})
	}
	return b.String()
}

// return the trailing whitespace of the input string
func trailingWhitespace(s string) string {
	var rev strings.Builder
	for i := len(s); i > 0; i-- {
		if !unicode.IsSpace(rune(s[i-1])) {
			return reverseString(rev.String())
		}
		rev.Write([]byte{s[i-1]})
	}
	return rev.String()
}

// reverse a string
func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, 0; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
