package http

import (
	"net/http"
)

// UserAgent is the default user agent for gopayment
const UserAgent = "GoPayment/2.0.0 (https://github.com/mohammadv184/gopayment)"

// transport is a custom RoundTripper to inject headers
type transport struct {
	inner http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	return t.inner.RoundTrip(req)
}

// NewClient returns a new http.client with a custom transport for GoPayment
func NewClient() *http.Client {
	return &http.Client{
		Transport: &transport{
			inner: http.DefaultTransport,
		},
	}
}
