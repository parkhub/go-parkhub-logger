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

// MARK: Re-used variables
var (
	mockHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(25 * time.Millisecond)
		(rw).Header().Set("Content-Type", "text/plain")
		(rw).WriteHeader(200)
		_, _ = (rw).Write([]byte("done"))
	})
	mockBody       = []byte(`{ "bodyProperty": "bodyValue" }`)
	rr             = httptest.NewRecorder()
	mockRequest, _ = http.NewRequest(
		http.MethodGet,
		"/test-path?testVar=testVal",
		bytes.NewBuffer(mockBody))
)

func TestMain(m *testing.M) {
	// setup
	mockRequest.Header.Set("Content-Type", "application/json")
	// run
	m.Run()
	// teardown
}
func noExtras(t *testing.T) {
	h := NewRequestLogger(RequestLoggerConfig{
		Headers: false,
		Params:  false,
		Body:    false,
	})
	h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
}
func headersAndParams(t *testing.T) {
	h := NewRequestLogger(RequestLoggerConfig{
		Headers: true,
		Params:  true,
		Body:    false,
	})
	h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
}
func headersAndBody(t *testing.T) {
	h := NewRequestLogger(RequestLoggerConfig{
		Headers: true,
		Params:  false,
		Body:    true,
	})
	h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
}
func timedOutRequest(t *testing.T) {
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
}
func cancelledRequest(t *testing.T) {
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
}
func handlerError(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("some other error"))
	errorRequest, _ := http.NewRequestWithContext(ctx, "GET", "/error-op", nil)
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Content-Type", "text/plain")
		(w).WriteHeader(500)
		_, _ = (w).Write([]byte("some error occurred"))
	})
	h := NewRequestLogger(RequestLoggerConfig{})
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
		h.Handle(mockHandler).ServeHTTP(rr, mockRequest)
	})

}
