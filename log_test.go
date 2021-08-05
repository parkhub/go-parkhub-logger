package log

import (
	"testing"
)

func TestSetupSingleLogger(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatPretty, false, true, []string{"test"})
	Infoln("This is an info statement.")
}

func TestSetupLocalLogger(t *testing.T) {
	SetupLocalLogger(LogLevelInfo)
	Infoln("This is an info statement.")
}

func TestSetupCloudLogger(t *testing.T) {
	SetupCloudLogger(LogLevelInfo, []string{"test", "tags"})
	Infoln("This is an info statement.")
}

func TestInfod(t *testing.T) {
	SetupLocalLogger(LogLevelDebug)

	type testStruct struct {
		Name string
		Kind string
	}

	test := &testStruct{
		Name: "Logan",
		Kind: "Log",
	}

	Infod("This is some data.", test)
}

func TestWhitespace(t *testing.T) {
	type testStruct struct {
		Name string
		Kind string
	}
	data := &testStruct{
		Name: "Logan",
		Kind: "Log",
	}
	message := `
	This is some data.
`

	t.Run("Pretty", func(t *testing.T) {
		SetupLocalLogger(LogLevelDebug)
		Infod(message, data)
	})

	t.Run("JSON", func(t *testing.T) {
		SetupCloudLogger(LogLevelDebug, []string{"logger", "test"})
		Infod(message, data)
	})
}

func TestDebugln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Debugln("This is a debug statement.")
}

func TestDebugf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Debugf("This is a debug statement %d.", 10000)
}

func TestInfoln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Infoln("This is an info statement.")
}

func TestInfof(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Infof("This is an info statement %d.", 10000)
}

func TestWarnln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Warnln("This is a warning.")
}

func TestWarnf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Warnf("This is a warning %d.", 10000)
}

func TestErrorln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Errorln("This is an error.")
}

func TestErrorf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Errorf("This is an error %d.", 10000)
}

func TestFatalln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Fatalln("This is fatal error.")
}

func TestFatalf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Fatalf("This is a fatal error %d.", 10000)
}
