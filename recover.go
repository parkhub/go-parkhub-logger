package log

import (
	"fmt"
	"runtime"
	"strings"
)

// Recover recovers from a panic to keep from crashing the application and
// logs the problem as an error. The error argument can be used to set an error
// to be returned by the calling function.
func Recover(in string, err *error) {
	if r := recover(); r != nil {
		var trace strings.Builder
		for i := 1; true; i++ {
			if pc, _, line, ok := runtime.Caller(i); ok {
				fn := runtime.FuncForPC(pc)
				trace.WriteString(fmt.Sprintf("\n\tfrom %s:%d", fn.Name(), line))
				continue
			}
			break
		}

		e, ok := r.(error)
		if ok {
			err = &e
			Errorf("%s panic: %v%s", in, e, trace.String())
			return
		}
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
		Errorf("%s panic: %v%s", in, r, trace.String())
	}
}