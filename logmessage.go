package log

import (
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/ttacon/chalk"
)

var leftRegexp = regexp.MustCompile(`^(\s)+`)
var rightRegexp = regexp.MustCompile(`(\s)+$`)

type logMessage struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Tags      []string    `json:"tags"`
	Message   string      `json:"message"`
	Metadata  interface{} `json:"metadata,omitempty"`
	File      string      `json:"file,omitempty"`

	format       Format `json:"-"`
	rawLevel     Level  `json:"-"`
	colorize     bool   `json:"-"`
	logCaller    bool   `json:"-"`
	trimmedLeft  string `json:"-"`
	trimmedRight string `json:"-"`
}

// newLogMessage creates a new log message.
func newLogMessage(format Format, colorize bool, logCaller bool, time logTime, level Level, tags []string, message string, data interface{}) *logMessage {
	modifiedMessage := message
	var trimmedLeft string
	var trimmedRight string

	left := leftRegexp.FindAllStringSubmatch(message, -1)
	if len(left) > 0 {
		trimmedLeft = left[0][0]
	}
	right := rightRegexp.FindAllStringSubmatch(message, -1)
	if len(right) > 0 {
		trimmedRight = right[0][0]
	}
	modifiedMessage = strings.TrimSpace(message)

	caller := ""
	if logCaller {
		var file string
		var line int
		var ok bool
		_, file, line, ok = runtime.Caller(4)
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

	formatedMessage := &logMessage{
		Timestamp:    time.String(),
		Level:        level.String(),
		Tags:         tags,
		Message:      modifiedMessage,
		Metadata:     data,
		File:         caller,
		format:       format,
		rawLevel:     level,
		colorize:     colorize,
		logCaller:    logCaller,
		trimmedLeft:  trimmedLeft,
		trimmedRight: trimmedRight,
	}

	return formatedMessage
}

// MARK: Methods

func (m logMessage) restoreWhitespace(output string) string {
	return m.trimmedLeft + output + m.trimmedRight
}

func (m logMessage) colorizeIfNeeded(output string) string {
	if !m.colorize || m.rawLevel == LogLevelDebug {
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
		m.trimmedLeft + m.Message + data + m.trimmedRight

	return m.colorizeIfNeeded(prettyMessage)
}
