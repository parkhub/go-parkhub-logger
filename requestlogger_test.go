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
var (
	mockHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(25 * time.Millisecond)
		(rw).Header().Set("Content-Type", "text/plain")
		(rw).WriteHeader(200)
		_, _ = (rw).Write([]byte("done"))
	})
	mockBody = []byte(`{ "bodyProperty": "bodyValue" }`)

	slowHandler = http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
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

	cancelHandler = http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
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

	errorHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Content-Type", "text/plain")
		(w).WriteHeader(500)
		_, _ = (w).Write([]byte("some error occurred"))
	})
)

func mockRequest() *http.Request {
	req, _ := http.NewRequest(
		http.MethodGet,
		"/test-path?testVar=testVal",
		bytes.NewBuffer(mockBody),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func noExtras(t *testing.T) {
	h := NewRequestLogger(RequestLoggerConfig{
		Headers: false,
		Params:  false,
		Body:    false,
	})
	rr := httptest.NewRecorder()
	h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
}

func headersAndParams(t *testing.T) {
	h := NewRequestLogger(RequestLoggerConfig{
		Headers: true,
		Params:  true,
		Body:    false,
	})
	rr := httptest.NewRecorder()
	h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
}

func headersAndBody(t *testing.T) {
	h := NewRequestLogger(RequestLoggerConfig{
		Headers: true,
		Params:  false,
		Body:    true,
	})
	rr := httptest.NewRecorder()
	h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
}

func timedOutRequest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
	h := NewRequestLogger(RequestLoggerConfig{})
	rr := httptest.NewRecorder()
	h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
}

func cancelledRequest(t *testing.T) {
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

	h := NewRequestLogger(RequestLoggerConfig{
		Tags:                  []string{"requests"},
		NormalLevel:           LogLevelTrace,
		DeadlineExceededLevel: LogLevelInfo,
		ContextCancelledLevel: LogLevelInfo,
		ContextErrorLevel:     LogLevelFatal,
	})
	rr := httptest.NewRecorder()
	t.Run("normal", func(t *testing.T) {
		h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
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

	t.Run("JSON Logger", func(t *testing.T) {
		SetupCloudLogger(LogLevelDebug, []string{"logger", "test"})

		t.Run("no extras", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: false,
				Params:  false,
				Body:    false,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
		})

		t.Run("headers and params", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: true,
				Params:  true,
				Body:    false,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
		})

		t.Run("headers and body", func(t *testing.T) {
			h := NewRequestLogger(RequestLoggerConfig{
				Headers: true,
				Params:  false,
				Body:    true,
			})
			h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
		})

		t.Run("timed-out rq", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 0)
			defer cancel()
			slowRequest, _ := http.NewRequestWithContext(ctx, "GET", "/slow-op", nil)
			h := NewRequestLogger(RequestLoggerConfig{})
			h.Handle(slowHandler).ServeHTTP(rr, slowRequest)
		})
	})
}

func handlerError(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("some other error"))
	errorRequest, _ := http.NewRequestWithContext(ctx, "GET", "/error-op", nil)
	h := NewRequestLogger(RequestLoggerConfig{})
	rr := httptest.NewRecorder()
	h.Handle(errorHandler).ServeHTTP(rr, errorRequest)
}

func TestRequestLogger_Handle(t *testing.T) {

	t.Run("String Logger", func(t *testing.T) {
		SetupLocalLogger(LogLevelDebug)
		t.Run("no extras", noExtras)
		t.Run("headers and params", headersAndParams)
		t.Run("headers and body", headersAndBody)
		t.Run("timed-out rq", timedOutRequest)
		t.Run("cancelled rq", cancelledRequest)
		t.Run("handler error", handlerError)
	})

	t.Run("JSON Logger", func(t *testing.T) {
		SetupCloudLogger(LogLevelDebug, []string{"logger", "test"})

		t.Run("no extras", noExtras)
		t.Run("headers and params", headersAndParams)
		t.Run("headers and body", headersAndBody)
		t.Run("timed-out rq", timedOutRequest)
		t.Run("cancelled rq", cancelledRequest)
		t.Run("handler error", handlerError)
	})

	t.Run("Sub-Logger", func(t *testing.T) {
		SetupLocalLogger(LogLevelDebug)
		sl := Sublogger("sub-logger")
		h := NewRequestLogger(RequestLoggerConfig{
			Logger: sl,
			Tags:   []string{"requests"},
		})
		rr := httptest.NewRecorder()
		h.Handle(mockHandler).ServeHTTP(rr, mockRequest())
	})

}
