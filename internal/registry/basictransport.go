package registry

import (
	"net/http"
)

type basicTransport struct {
	next     http.RoundTripper
	host     string
	username string
	password string
}

func NewBasicTransport(host, username, password string, next http.RoundTripper) *basicTransport {
	return &basicTransport{
		next:     next,
		host:     host,
		username: username,
		password: password,
	}
}

func (t *basicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.host == req.URL.Host {
		if t.username != "" || t.password != "" {
			req.SetBasicAuth(t.username, t.password)
		}
	}
	resp, err := t.next.RoundTrip(req)
	return resp, err
}
