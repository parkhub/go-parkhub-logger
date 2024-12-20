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
	client                *http.Client
	logger                Logger
	logHeaders            bool
	logParams             bool
	logBody               bool
	normalLevel           Level
	deadlineExceededLevel Level
	contextCancelledLevel Level
	contextErrorLevel     Level
}

// RequestLoggerConfig defines options for which details should be logged
type RequestLoggerConfig struct {
	Logger  Logger
	Client  *http.Client
	Headers bool
	Params  bool
	GraphQL bool
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
	graphql      string
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
			rl.logger.Errord("error creating request log for "+r.URL.String()+":", err)
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
	l := config.Logger
	if l == nil {
		l = LoggerSingleton
	}
	sl := l.Sublogger(config.Tags...)
	sl.(*sublogger).skipOffset += 3
	normal, deadline, cancelled, ctxErr := LogLevelDebug, LogLevelWarn, LogLevelWarn, LogLevelError
	if config.NormalLevel != logLevelUnset {
		normal = config.NormalLevel
	}
	if config.DeadlineExceededLevel != logLevelUnset {
		deadline = config.DeadlineExceededLevel
	}
	if config.ContextCancelledLevel != logLevelUnset {
		cancelled = config.ContextCancelledLevel
	}
	if config.ContextErrorLevel != logLevelUnset {
		ctxErr = config.ContextErrorLevel
	}
	return &RequestLogger{
		logger:                sl,
		client:                config.Client,
		logHeaders:            config.Headers,
		logParams:             config.Params,
		logBody:               config.Body,
		normalLevel:           normal,
		deadlineExceededLevel: deadline,
		contextCancelledLevel: cancelled,
		contextErrorLevel:     ctxErr,
	}
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

	rl.logger.Logd(ll, log.label(), log)
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
	if opts.GraphQL {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return log, fmt.Errorf("error reading request body: %s", err), http.StatusBadRequest
		}
		r.Body.Close()

		var bodyJson map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &bodyJson); err == nil {
			opName, opOk := bodyJson["operationName"].(string)
			q, qOk := bodyJson["query"].(string)
			if opOk && qOk && opName != "" && q != "" {
				// Only log properly formatted graphql requests
				// (json with operationName and query)
				if strings.HasPrefix(q, "query") {
					log.graphql = fmt.Sprintf("q %s", opName)
				} else if strings.HasPrefix(q, "mutation") {
					log.graphql = fmt.Sprintf("m %s", opName)
				}
			}
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
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
