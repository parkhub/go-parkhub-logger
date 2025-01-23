package log

import (
	"net/http"
	"time"
)

// MARK: Public Functions

func (rl *RequestLogger) Do(req *http.Request) (res *http.Response, err error) {
	return rl.RoundTrip(req)
}

func (rl *RequestLogger) RoundTrip(req *http.Request) (res *http.Response, err error) {
	client := rl.client
	if client == nil {
		client = http.DefaultClient
	}
	log, err, _ := makeLog(req, RequestLoggerConfig{
		Headers: rl.logHeaders,
		Params:  rl.logParams,
		Body:    rl.logBody,
	})
	if err != nil {
		rl.Logger.Errord("error creating request log for "+req.URL.String()+":", err)
		return rl.client.Do(req)
	}

	start := time.Now()
	res, err = rl.client.Do(req)
	log.latency = time.Since(start)
	log.contextError = req.Context().Err()
	rl.log(log)

	return
}

// MARK: Interface Checks
var (
	_ http.RoundTripper = (*RequestLogger)(nil)
)
