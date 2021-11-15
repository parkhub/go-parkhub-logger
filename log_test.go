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

func TestTrace(t *testing.T) {
	SetupLogger(LogLevelTrace, LogFormatJSON, false, true, []string{"test", "tags"})

	t.Run("Traceln", func(t *testing.T) {
		Traceln("This is a trace ln statement.")
	})

	t.Run("Tracef", func(t *testing.T) {
		Tracef("This is a trace f statement %d.", 10000)
	})

	t.Run("Traced", func(t *testing.T) {
		Traced("This is a trace d statement.", 10000)
	})
}

func TestDebug(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})

	t.Run("Debugln", func(t *testing.T) {
		Debugln("This is a debug ln statement.")
	})

	t.Run("Debugf", func(t *testing.T) {
		Debugf("This is a debug f statement %d.", 10000)
	})

	t.Run("Debugd", func(t *testing.T) {
		Debugd("This is a debug d statement.", 10000)
	})
}

func TestInfo(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})

	t.Run("Infoln", func(t *testing.T) {
		Infoln("This is an info ln statement.")
	})

	t.Run("Infof", func(t *testing.T) {
		Infof("This is an info f statement %d.", 10000)
	})

	t.Run("Infod", func(t *testing.T) {
		Infod("This is an info f statement.", 10000)
	})
}

func TestWarn(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})

	t.Run("Warnln", func(t *testing.T) {
		Warnln("This is a warning ln.")
	})

	t.Run("Warnf", func(t *testing.T) {
		Warnf("This is a warning f %d.", 10000)
	})

	t.Run("Warnd", func(t *testing.T) {
		Warnd("This is a warning d.", 10000)
	})
}

func TestError(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})

	t.Run("Errorln", func(t *testing.T) {
		Errorln("This is an error ln.")
	})

	t.Run("Errorf", func(t *testing.T) {
		Errorf("This is an error f %d.", 10000)
	})

	t.Run("Errord", func(t *testing.T) {
		Errord("This is an error d.", 10000)
	})
}

func TestFatal(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatJSON, false, true, []string{"test", "tags"})

	t.Run("Fatalln", func(t *testing.T) {
		Fatalln("This is a fatal ln.")
	})

	t.Run("Fatalf", func(t *testing.T) {
		Fatalf("This is a fatal f %d.", 10000)
	})

	t.Run("Fatald", func(t *testing.T) {
		Fatald("This is a fatal d.", 10000)
	})
}
