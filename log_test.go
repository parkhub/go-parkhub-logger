package log

import (
	"testing"
)

func TestSetupSingleLogger(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatPretty, false, []string{"test"})
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
	Debugln("This is a debug statement.")
}

func TestDebugf(t *testing.T) {
	Debugf("This is a debug statement %d.\n", 10000)
}

func TestInfoln(t *testing.T) {
	Infoln("This is an info statement.")
}

func TestInfof(t *testing.T) {
	Infof("This is an info statement %d.\n", 10000)
}

func TestWarnln(t *testing.T) {
	Warnln("This is a warning.")
}

func TestWarnf(t *testing.T) {
	Warnf("This is a warning %d.\n", 10000)
}

func TestErrorln(t *testing.T) {
	Errorln("This is an error.")
}

func TestErrorf(t *testing.T) {
	Errorf("This is an error %d.\n", 10000)
}

func TestFatalln(t *testing.T) {
	Fatalln("This is fatal error.")
}

func TestFatalf(t *testing.T) {
	Fatalf("This is a fatal error %d.\n", 10000)
}
