package log

import (
	"net/http"
	"time"
)

// MARK: Public Functions

func (rl *RequestLogger) RoundTripper(client *http.Client) http.RoundTripper {
	var c *http.Client
	if client != nil {
		c = client
	} else if rl.client != nil {
		c = rl.client
	} else {
		c = http.DefaultClient
	}
	return roundTripper{
		RequestLogger: rl,
		Client:        c,
	}
}

func (rt roundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	log, err, _ := makeLog(req, RequestLoggerConfig{
		Headers: rt.RequestLogger.logHeaders,
		Params:  rt.RequestLogger.logParams,
		Body:    rt.RequestLogger.logBody,
	})
	if err != nil {
		rt.logger.Errord("error creating request log for "+req.URL.String()+":", err)
		return rt.Client.Do(req)
	}

	start := time.Now()
	res, err = rt.Client.Do(req)
	log.latency = time.Since(start)
	log.contextError = req.Context().Err()
	rt.log(log)

	return
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
		rl.logger.Errord("error creating request log for "+req.URL.String()+":", err)
		return rl.client.Do(req)
	}

	start := time.Now()
	res, err = rl.client.Do(req)
	log.latency = time.Since(start)
	log.contextError = req.Context().Err()
	rl.log(log)

	return
}

// MARK: Private Types

type roundTripper struct {
	*RequestLogger
	*http.Client
}
