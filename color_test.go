package log

import "testing"

func TestColor(t *testing.T) {
	SetupLogger(LogLevelTrace, LogFormatPretty, true, false, []string{"test"})

	Traceln("Trace")
	Debugln("Debug")
	Infoln("Info")
	Warnln("Warn")
	Errorln("Error")

	t.Skip("Skipping fatal check because it exits")
	Fatalln("Fatal")
}
