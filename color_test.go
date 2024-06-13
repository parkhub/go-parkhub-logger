package log

import (
	"testing"
)

func TestColor(t *testing.T) {
	LoggerSingleton = &logger{
		rawLevel:       LogLevelTrace,
		format:         LogFormatPretty,
		tags:           []string{"Environment", "Platform", "Application"},
		colorizeOutput: true,
		logCaller:      true,
		exitFunc:       func() { t.Log("> os.Exit(1)") },
	}

	Logln(10, "Default")

	Traceln("This is a trace message.")
	Debugln("This is a debug message.")
	Infoln("This is an info message.")
	Warnln("This is a warn message.")
	Errorln("This is an error message.")
	Fatalln("This is a fatal error message.")
}
