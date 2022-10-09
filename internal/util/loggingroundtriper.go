package util

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
)

type LoggingRoundTripper struct {
	next http.RoundTripper
}

func NewLoggingRoundTripper(next http.RoundTripper) http.RoundTripper {
	return &LoggingRoundTripper{
		next: next,
	}
}

func (l LoggingRoundTripper) logRequest(req *http.Request) {
	log.Debug().Msgf(">>> %v", req.URL)
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Trace().Err(err).Msgf("failed to dump request")
	}
	log.Trace().Msg(string(requestDump))
}

func (l LoggingRoundTripper) logResponse(res *http.Response, err error) {
	if err != nil {
		log.Err(err)
	} else {
		log.Debug().Msgf("<<< %v", res.Request.URL)
		responseDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Debug().Err(err).Msg("failed to dump response")
		}
		log.Trace().Msgf(string(responseDump))
	}
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	l.logRequest(req)
	res, err = l.next.RoundTrip(req)
	l.logResponse(res, err)
	return
}
