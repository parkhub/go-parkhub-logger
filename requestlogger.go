package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MARK: Types

// RequestLogger is a configured struct with an HTTP Handler method for logging
// HTTP requests
type RequestLogger struct {
	Logger
	client                *http.Client
	tags                  []string
	logHeaders            bool
	logParams             bool
	logBody               bool
	normalLevel           Level
	deadlineExceededLevel Level
	contextCancelledLevel Level
	contextErrorLevel     Level
}

// RequestLoggerConfig defines options for which details should be logged.
// Setting a LogLevel to 0 will prevent it from being logged.
type RequestLoggerConfig struct {
	Logger  Logger
	Client  *http.Client
	Headers bool
	Params  bool
	Body    bool
	Tags    []string

	// NormalLevel is the log level to use for requests without context errors
	NormalLevel Level

	// DeadlineExceededLevel is the log level to use for requests that time out
	DeadlineExceededLevel Level

	// ContextCancelledLevel is the log level to use for requests that are
	// cancelled for reasons other than a timeout
	ContextCancelledLevel Level

	// ContextErrorLevel is the log level to use for requests that have an other
	// context error
	ContextErrorLevel Level
}

// requestLog stores the request data for logging
type requestLog struct {
	headers      map[string][]string
	method       string
	path         string
	params       map[string][]string
	body         string
	latency      time.Duration
	contextError error
}

// MARK: Public Methods

// Handle logs incoming HTTP requests, calls the next handler, and logs uncaught
// errors in the handler chain
func (rl *RequestLogger) Handle(next http.Handler) http.Handler {
	opts := RequestLoggerConfig{
		Headers: rl.logHeaders,
		Params:  rl.logParams,
		Body:    rl.logBody,
	}
	logChan := make(chan requestLog)

	go func() {
		log := <-logChan
		close(logChan)
		rl.log(log)
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, err, statusCode := makeLog(r, opts)
		if err != nil {
			rl.Logger.Errord("error creating request log for "+r.URL.String()+":", err)
			http.Error(w, err.Error(), statusCode)
			return
		}

		start := time.Now().UTC()
		next.ServeHTTP(w, r)
		end := time.Now().UTC()
		log.latency = end.Sub(start)
		log.contextError = r.Context().Err()
		logChan <- log
	})
}

// String returns the requestLog as a formatted string
func (rl requestLog) String() string {
	var headerStr, paramStr, bodyStr string

	// format headers
	if i, l := 0, len(rl.headers); l > 0 {
		headers := make([]string, l)
		for k, v := range rl.headers {
			headers[i] = fmt.Sprintf("%s: [%s]", k, strings.Join(v, ", "))
			i++
		}
		headerStr = "\nHeaders: " + strings.Join(headers, "; ")
	}

	// format query params
	if i, l := 0, len(rl.params); l > 0 {
		params := make([]string, l)
		for k, v := range rl.params {
			params[i] = fmt.Sprintf("%s: [%s]", k, strings.Join(v, ", "))
			i++
		}
		paramStr = "\nParams: " + strings.Join(params, "; ")
	}

	// format requestBody
	if rl.body != "" {
		bodyStr = "\nBody: " + rl.body
	}

	return headerStr + paramStr + bodyStr
}

// MarshalJSON returns the requestLog as a JSON object
func (rl requestLog) MarshalJSON() ([]byte, error) {
	obj := map[string]interface{}{
		"method":  rl.method,
		"path":    rl.path,
		"latency": rl.latency,
	}
	if len(rl.headers) > 0 {
		obj["headers"] = rl.headers
	}
	if err := rl.contextError; err != nil {
		obj["cancelReason"] = err.Error()
	}
	if len(rl.params) > 0 {
		obj["params"] = rl.params
	}
	if len(rl.body) > 0 {
		obj["body"] = rl.body
	}

	return json.Marshal(obj)
}

// MARK: Public Functions

// NewRequestLogger returns a configured RequestLogger
func NewRequestLogger(config RequestLoggerConfig) *RequestLogger {
	config = defaultRLConfig.apply(config)
	l := config.Logger
	sl := l.Sublogger(config.Tags...)
	sl.(*sublogger).skipOffset += 3
	return &RequestLogger{
		Logger:                sl,
		client:                config.Client,
		tags:                  config.Tags,
		logHeaders:            config.Headers,
		logParams:             config.Params,
		logBody:               config.Body,
		normalLevel:           config.NormalLevel,
		deadlineExceededLevel: config.DeadlineExceededLevel,
		contextCancelledLevel: config.ContextCancelledLevel,
		contextErrorLevel:     config.ContextErrorLevel,
	}
}

// LegacyDefaults returns a new RequestLoggerConfig with the historical
// default values applied to unset values
//
// DEPRECATED
func (cfg RequestLoggerConfig) LegacyDefaults() RequestLoggerConfig {
	return RequestLoggerConfig{
		Logger:                LoggerSingleton,
		Client:                http.DefaultClient,
		Headers:               false,
		Params:                false,
		Body:                  false,
		Tags:                  nil,
		NormalLevel:           LogLevelDebug,
		DeadlineExceededLevel: LogLevelWarn,
		ContextCancelledLevel: LogLevelWarn,
		ContextErrorLevel:     LogLevelError,
	}.apply(defaultRLConfig)
}

// MARK: Private Methods

// log logs a requestLog with the RequestLogger's Logger
func (rl *RequestLogger) log(log requestLog) {
	var ll Level
	switch {
	case errors.Is(log.contextError, context.DeadlineExceeded):
		ll = rl.deadlineExceededLevel
	case errors.Is(log.contextError, context.Canceled):
		ll = rl.contextCancelledLevel
	case log.contextError != nil:
		ll = rl.contextErrorLevel
	default:
		ll = rl.normalLevel
	}

	rl.Logger.Logd(ll, log.label(), log)
}

func (rl *RequestLogger) config() RequestLoggerConfig {
	return RequestLoggerConfig{
		Logger:                rl.Logger,
		Client:                rl.client,
		Headers:               rl.logHeaders,
		Params:                rl.logParams,
		Body:                  rl.logBody,
		Tags:                  rl.tags,
		NormalLevel:           rl.normalLevel,
		DeadlineExceededLevel: rl.deadlineExceededLevel,
		ContextCancelledLevel: rl.contextCancelledLevel,
		ContextErrorLevel:     rl.contextErrorLevel,
	}
}

func makeLog(r *http.Request, opts RequestLoggerConfig) (requestLog, error, int) {
	log := requestLog{
		method: r.Method,
		path:   r.URL.Path,
	}

	if opts.Headers {
		log.headers = r.Header
	}
	if opts.Params {
		log.params = r.URL.Query()
	}
	if opts.Body {
		buf, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			return log, errors.New("Failed to read request body: " + bodyErr.Error()), http.StatusBadRequest
		}
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		bodyBytes, _ := io.ReadAll(rdr1)

		log.body = string(bodyBytes)
		var isJSON bool
		for _, v := range r.Header["Content-Type"] {
			if v == "application/json" {
				isJSON = true
				break
			}
		}
		if isJSON {
			rawJSON := json.RawMessage(bodyBytes)
			if jsonBytes, err := json.MarshalIndent(rawJSON, "", "  "); err == nil {
				log.body = string(jsonBytes)
			}
		}

		r.Body = rdr2
	}

	return log, nil, 0
}

// label returns the message to print to the log
func (rl requestLog) label() string {
	var label string
	switch {
	case errors.Is(rl.contextError, context.DeadlineExceeded):
		label = fmt.Sprintf("%s %s: %dms (DEADLINE EXCEEDED)", rl.method, rl.path, rl.latency)
	case errors.Is(rl.contextError, context.Canceled):
		label = fmt.Sprintf("%s %s: %dms (CANCELLED)", rl.method, rl.path, rl.latency)
	case rl.contextError != nil:
		label = fmt.Sprintf("%s %s: %dms (%s)", rl.method, rl.path, rl.latency, rl.contextError)
	default:
		label = fmt.Sprintf("%s %s: %dms", rl.method, rl.path, rl.latency)
	}
	return label
}

// apply returns a new RequestLoggerConfig from config, using the receiver's
// values where config's values are unset.
// The Headers, Params, and Body fields will always be cfg2's values, even if
// not explicitly set.
func (cfg RequestLoggerConfig) apply(config RequestLoggerConfig) RequestLoggerConfig {
	if config.Logger == nil {
		config.Logger = cfg.Logger
	}
	if config.Client == nil {
		config.Client = cfg.Client
	}
	if config.NormalLevel == logLevelUnset {
		config.NormalLevel = cfg.NormalLevel
	}
	if config.DeadlineExceededLevel == logLevelUnset {
		config.DeadlineExceededLevel = cfg.DeadlineExceededLevel
	}
	if config.ContextCancelledLevel == logLevelUnset {
		config.ContextCancelledLevel = cfg.ContextCancelledLevel
	}
	if config.ContextErrorLevel == logLevelUnset {
		config.ContextErrorLevel = cfg.ContextErrorLevel
	}

	return config
}

// MARK: Private Variables

var defaultRLConfig = RequestLoggerConfig{
	Logger:                LoggerSingleton,
	Client:                nil,
	Headers:               false,
	Params:                false,
	Body:                  false,
	Tags:                  nil,
	NormalLevel:           LogLevelDebug,
	DeadlineExceededLevel: LogLevelWarn,
	ContextCancelledLevel: LogLevelWarn,
	ContextErrorLevel:     LogLevelError,
}

// MARK: Interface Checks
var (
	_ Logger       = (*RequestLogger)(nil)
	_ http.Handler = new(RequestLogger).Handle(http.HandlerFunc(nil))
)
