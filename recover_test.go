package log

import (
	"testing"
)

func causePanic() {
	panic("at the disco!")
}

func recoverPanic() (err error) {
	defer Recover("recoverPanic", &err)
	causePanic()
	return nil
}

func TestRecover(t *testing.T) {
	SetupLogger(LogLevelDebug, LogFormatPretty, false, false, nil)

	t.Run("No Error", func(t *testing.T) {
		defer Recover("NoError", nil)
	})

	t.Run("Catch Error", func(t *testing.T) {
		if err := recoverPanic(); err == nil {
			t.Error("error not set")
		}
	})

	t.Run("Index Out Of Range", func(t *testing.T) {
		defer Recover("Test Index Out Of Range", nil)
		_ = []struct{}{}[0]
	})

	t.Run("Nil Pointer Dereference", func(t *testing.T) {
		defer Recover("Test Nil Pointer Dereference", nil)
		var ref *struct{}
		_ = *ref
	})

	t.Run("Non-Error Panic", func(t *testing.T) {
		defer Recover("Test Non-Error Panic", nil)
		panic("at the disco!")
	})
}
