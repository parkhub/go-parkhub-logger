package log

import (
	"fmt"
	"testing"
)

func TestSublogger(t *testing.T) {
	LoggerSingleton = &logger{
		rawLevel:       LogLevelDebug,
		format:         LogFormatPretty,
		tags:           []string{"Environment", "Platform", "Application"},
		colorizeOutput: false,
		logCaller:      true,
		exitFunc:       func() { fmt.Println("> os.Exit(1)") },
	}

	sl := Sublogger("Function")
	sl2 := sl.Sublogger("Sub-Function")
	Debugln("Singleton Logger")
	sl.Debugln("sublogger")
	sl2.Debugln("sub-sub-logger")
}
