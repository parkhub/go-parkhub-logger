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
	sl2 := sl.Sublogger("Sub-Function")
	Debugln("Singleton Logger")
	sl.Debugln("sublogger")
	sl2.Debugln("sub-sub-logger")
}
