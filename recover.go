package log

// Recover recovers from a panic to keep from crashing the application and
// logs the problem as an error
func Recover(in string) {
	if r := recover(); r != nil {
		e, ok := r.(error)
		if ok {
			Errord(in+" panic:", e)
			return
		}
		Errorf("%s panic: %v", in, r)
	}
}
