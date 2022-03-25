package log

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MARK: Test Types

// errContext is a context that immediately closes and returns an error that's
// neither context.Canceled nor context.DeadlineExceeded
type errContext struct {
	done chan struct{}
}

func (ec errContext) Deadline() (deadline time.Time, ok bool) { return }
func (ec errContext) Done() <-chan struct{} {
	defer func() {
		if ec.done != nil {
			ec.done <- struct{}{}
		}
	}()
	return ec.done
}
func (ec errContext) Err() error {
	return errors.New("unusual context error")
}
func (ec errContext) Value(key interface{}) interface{} { return nil }

// MARK: Re-used variables
var mockHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	time.Sleep(25 * time.Millisecond)
	(rw).Header().Set("Content-Type", "text/plain")
	(rw).WriteHeader(200)
	_, _ = (rw).Write([]byte("done"))
})
var mockBody = []byte(`{ "bodyProperty": "bodyValue" }`)
var slowHandler = http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
	res := make(chan struct{})
	go func() {
		time.Sleep(25 * time.Millisecond)
		res <- struct{}{}
	}()

	select {
	case <-rq.Context().Done():
		(rw).Header().Set("Content-Type", "text/plain")
		(rw).WriteHeader(500)
		_, _ = (rw).Write([]byte("timeout"))
	case <-res:
		(rw).Header().Set("Content-Type", "text/plain")
		(rw).WriteHeader(200)
		_, _ = (rw).Write([]byte("done"))
	}
})
var cancelHandler = http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
	res := make(chan struct{})
	go func() {
		time.Sleep(25 * time.Millisecond)
		res <- struct{}{}
	}()

	select {
	case <-rq.Context().Done():
		(rw).Header().Set("Content-Type", "text/plain")
		(rw).WriteHeader(500)
		_, _ = (rw).Write([]byte("canceled"))
	case <-res:
		(rw).Header().Set("Content-Type", "text/plain")
		(rw).WriteHeader(200)
		_, _ = (rw).Write([]byte("done"))
	}
})

func TestRequestLogger_Handle(t *testing.T) {
	rr := httptest.NewRecorder()
	var mockRequest, _ = http.NewRequest(
		http.MethodGet,
		"/test-path?testVar=testVal",
		bytes.NewBuffer(mockBody))
	mockRequest.Header.Set("Content-Type", "application/json")

	t.Run("String Logger", func(t *testing.T) {
		SetupLocalLogger(LogLevelTrace)

		t.Run("no extras", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
		})

		t.Run("headers and params", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: true,
				Params:  true,
				Body:    false,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
		})

		t.Run("headers and body", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: true,
				Params:  false,
				Body:    true,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
		})

		t.Run("timed-out rq", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()
			slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)

			h := NewRequestLogger(RequestLoggerConfig{})
			h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
		})

		t.Run("cancelled rq", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)

			h := NewRequestLogger(RequestLoggerConfig{})
			h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
		})

		t.Run("custom levels", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Tags:                  []string{"requests"},
				NormalLevel:           LogLevelTrace,
				DeadlineExceededLevel: LogLevelInfo,
				ContextCancelledLevel: LogLevelInfo,
				ContextErrorLevel:     LogLevelFatal,
			})
			t.Run("normal", func(t *testing.T) {
				h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
			})
			t.Run("timeout", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 0)
				defer cancel()
				slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
				h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
			})
			t.Run("cancel", func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
				h.Handle(cancelHandler).ServeHTTP(rr, slowRequest)
			})
			t.Run("other error", func(t *testing.T) {
				ctx := errContext{}
				slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
				h.Handle(cancelHandler).ServeHTTP(rr, slowRequest)
			})
		})
	})

	t.Run("JSON Logger", func(t *testing.T) {
		SetupCloudLogger(LogLevelDebug, []string{"logger", "test"})

		t.Run("no extras", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: false,
				Params:  false,
				Body:    false,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
		})

		t.Run("headers and params", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: true,
				Params:  true,
				Body:    false,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
		})

		t.Run("headers and body", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: true,
				Params:  false,
				Body:    true,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
		})

		t.Run("timed-out rq", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()
			slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
			slowHandler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
				res := make(chan struct{})
				go func() {
					time.Sleep(25 * time.Millisecond)
					res <- struct{}{}
				}()

				select {
				case <-rq.Context().Done():
					(rw).Header().Set("Content-Type", "text/plain")
					(rw).WriteHeader(500)
					_, _ = (rw).Write([]byte("timeout"))
				case <-res:
					(rw).Header().Set("Content-Type", "text/plain")
					(rw).WriteHeader(200)
					_, _ = (rw).Write([]byte("done"))
				}
			})

			h := NewRequestLogger(RequestLoggerConfig{})
			h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
		})

		t.Run("cancelled rq", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
			slowHandler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
				res := make(chan struct{})
				go func() {
					time.Sleep(25 * time.Millisecond)
					res <- struct{}{}
				}()

				select {
				case <-rq.Context().Done():
					(rw).Header().Set("Content-Type", "text/plain")
					(rw).WriteHeader(500)
					_, _ = (rw).Write([]byte("canceled"))
				case <-res:
					(rw).Header().Set("Content-Type", "text/plain")
					(rw).WriteHeader(200)
					_, _ = (rw).Write([]byte("done"))
				}
			})

			h := NewRequestLogger(RequestLoggerConfig{})
			h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
		})
	})

	t.Run("Sub-Logger", func(t *testing.T) {
		SetupLocalLogger(LogLevelDebug)
		sl := Sublogger("sub-logger")
		h := NewRequestLogger(RequestLoggerConfig{
			Logger: sl,
			Tags:   []string{"requests"},
		})
		h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
	})

}
