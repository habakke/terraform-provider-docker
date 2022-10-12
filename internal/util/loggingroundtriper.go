package util

import (
	"context"
	"net/http"
	"net/http/httputil"
)

type Logger interface {
	Tracef(ctx context.Context, msg string, additionalFields ...interface{})
	Debugf(ctx context.Context, msg string, additionalFields ...interface{})
	Warnf(ctx context.Context, msg string, additionalFields ...interface{})
	Errorf(ctx context.Context, msg string, additionalFields ...interface{})
	Infof(ctx context.Context, msg string, additionalFields ...interface{})
}

type LoggingRoundTripper struct {
	next   http.RoundTripper
	logger Logger
	ctx    context.Context
}

func NewLoggingRoundTripper(ctx context.Context, next http.RoundTripper, logger Logger) http.RoundTripper {
	return &LoggingRoundTripper{
		next:   next,
		logger: logger,
		ctx:    ctx,
	}
}

func (l LoggingRoundTripper) logRequest(req *http.Request) {
	l.logger.Debugf(l.ctx, ">>> %v", req.URL)
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		l.logger.Errorf(l.ctx, "failed to dump request")
	}
	l.logger.Tracef(l.ctx, string(requestDump))
}

func (l LoggingRoundTripper) logResponse(res *http.Response, err error) {
	if err != nil {
		l.logger.Errorf(l.ctx, err.Error())
	} else {
		l.logger.Debugf(l.ctx, "<<< %v", res.Request.URL)
		responseDump, err := httputil.DumpResponse(res, true)
		if err != nil {
			l.logger.Debugf(l.ctx, "failed to dump response: %v", err.Error())
		}
		l.logger.Tracef(l.ctx, string(responseDump))
	}
}

func (l LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	l.logRequest(req)
	res, err = l.next.RoundTrip(req)
	l.logResponse(res, err)
	return
}
