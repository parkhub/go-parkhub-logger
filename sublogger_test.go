package log

import "testing"

func TestSublogger(t *testing.T) {
	SetupLogger(
		LogLevelDebug,
		LogFormatPretty,
		false,
		false,
		[]string{"Environment", "Platform", "Application"},
	)
	sl := Sublogger("Function")
	Debugln("Singleton Logger")
	sl.Debugln("sublogger")
}
