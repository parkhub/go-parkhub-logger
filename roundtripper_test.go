package log

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name          string
		request       *http.Request
		errorExpected bool
	}{
		{
			name: "valid request",
			request: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "https",
					Host:   "nowhere.net",
					Path:   "/testpath",
				},
			},
			errorExpected: false,
		},
		{
			name: "error in request",
			request: &http.Request{
				Method: "\\",
				URL:    &url.URL{Path: "test-path"},
			},
			errorExpected: true,
		},
	}

	SetupLocalLogger(LogLevelTrace)
	client := http.DefaultClient
	rl := NewRequestLogger(RequestLoggerConfig{
		Tags:                  []string{"requests"},
		Client:                client,
		NormalLevel:           LogLevelTrace,
		DeadlineExceededLevel: LogLevelInfo,
		ContextCancelledLevel: LogLevelInfo,
		ContextErrorLevel:     LogLevelFatal,
	})
	rt := rl.RoundTripper(client)
	responseRecorder := httptest.NewRecorder()

	for _, tt := range tests {
		for _, rt := range []http.RoundTripper{rt, rl} {
			t.Run(tt.name, func(t *testing.T) {
				_, err := rt.RoundTrip(tt.request)
				if (err != nil) != tt.errorExpected {
					t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.errorExpected)
					return
				}

				err = responseRecorder.Result().Body.Close()
				responseData, err := io.ReadAll(responseRecorder.Result().Body)

				if err != nil {
					t.Failed()
				}
				responseString := string(responseData)

				if responseString != "" {
					t.Errorf("Expected empty response body, got %s", responseString)
				}
			})
		}
	}
}
