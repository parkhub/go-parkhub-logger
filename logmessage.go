package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ttacon/chalk"
)

type logMessage struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Tags      []string    `json:"tags"`
	Message   string      `json:"message"`
	Metadata  interface{} `json:"metadata,omitempty"`

	format         Format `json:"-"`
	rawLevel       Level  `json:"-"`
	colorize       bool   `json:"-"`
	trimmedNewline bool   `json:"-"`
}

// newLogMessage creates a new log message.
func newLogMessage(format Format, colorize bool, time logTime, level Level, tags []string, message string, data interface{}) *logMessage {
	modifiedMessage := message
	hasSuffix := false

	if strings.HasSuffix(message, "\n") {
		// Strip newline from end of message
		hasSuffix = true
		modifiedMessage = message[:len(message)-1]
	}

	formatedMessage := &logMessage{
		Timestamp:      time.String(),
		Level:          level.String(),
		Tags:           tags,
		Message:        modifiedMessage,
		Metadata:       data,
		format:         format,
		rawLevel:       level,
		colorize:       colorize,
		trimmedNewline: hasSuffix,
	}

	return formatedMessage
}

// MARK: Methods

func (m logMessage) restoreNewline(output string) string {
	tail := ""
	if m.trimmedNewline {
		tail = "\n"
	}

	return output + tail
}

func (m logMessage) colorizeIfNeeded(output string) string {
	if !m.colorize || m.rawLevel == LogLevelDebug {
		return output
	}

	return fmt.Sprintf("%s%s%s", m.rawLevel.color(), output, chalk.ResetColor)
}

func (m logMessage) jsonString() string {
	s, _ := json.Marshal(m)
	return m.restoreNewline(m.colorizeIfNeeded(string(s)))
}

// MARK: String interface methods

func (m logMessage) String() string {
	if m.format == LogFormatJSON {
		return m.jsonString()
	}

	prettyMessage := fmt.Sprintf("%s [%s] ", m.Timestamp, m.Level)

	if len(m.Tags) > 0 {
		prettyMessage += fmt.Sprintf("(%s) ", strings.Join(m.Tags, ","))
	}

	prettyMessage += m.Message

	if m.Metadata != nil {
		prettyMessage += fmt.Sprintf(" %+v", m.Metadata)
	}

	return m.restoreNewline(m.colorizeIfNeeded(prettyMessage))
}
