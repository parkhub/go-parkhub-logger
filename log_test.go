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

	Infodln("This is some data.", test)
}

func TestDebugln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Debugln("This is a debug statement.")
}

func TestDebugf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Debugf("This is a debug statement %d.\n", 10000)
}

func TestInfoln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Infoln("This is an info statement.")
}

func TestInfof(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Infof("This is an info statement %d.\n", 10000)
}

func TestWarnln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Warnln("This is a warning.")
}

func TestWarnf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Warnf("This is a warning %d.\n", 10000)
}

func TestErrorln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Errorln("This is an error.")
}

func TestErrorf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Errorf("This is an error %d.\n", 10000)
}

func TestFatalln(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Fatalln("This is fatal error.")
}

func TestFatalf(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})
	Fatalf("This is a fatal error %d.\n", 10000)
}
