package log

import (
	"testing"
	"time"
)

func TestSetupSingleLogger(t *testing.T) {
	SetupLogger(0, []string{"test"}, false, Pretty)

	Infoln("This is an info statement.")
}

func TestInfod(t *testing.T) {
	SetupLogger(0, []string{"test"}, false, Pretty)

	type testStruct struct {
		Name string
		Kind string
	}

	test := &testStruct{
		Name: "Logan",
		Kind: "Log",
	}

	Infod("This is some data", test)

	time.Sleep(3 * time.Second)
}

func TestDebugln(t *testing.T) {
	Debugln("This is a debug statement.")
}

func TestDebugf(t *testing.T) {
	Debugf("This is a debug statement %d.", 10000)
}

func TestInfoln(t *testing.T) {
	Infoln("This is an info statement.")
}

func TestInfof(t *testing.T) {
	Infof("This is an info statement %d.", 10000)
}

func TestWarnln(t *testing.T) {
	Warnln("This is a warning.")
}

func TestWarnf(t *testing.T) {
	Warnf("This is a warning %d.", 10000)
}

func TestErrorln(t *testing.T) {
	Errorln("This is an error.")
}

func TestErrorf(t *testing.T) {
	Errorf("This is an error %d.", 10000)
}

func TestFatalln(t *testing.T) {
	Fatalln("This is an error.")
}

func TestFatalf(t *testing.T) {
	Fatalf("This is an error %d.", 10000)
}
