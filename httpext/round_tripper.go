package httpext

import (
	"net/http"
	"time"
)

// RoundTripper is a custom HTTP round tripper that implements the http.RoundTripper interface
// this Roundtripper will be used when transport options are set in the customClient
type RoundTripper struct {
	base http.RoundTripper
}

func NewRoundTripper(maxIdleConnsPerHost int, idleConnTimeout time.Duration) *RoundTripper {
	return &RoundTripper{
		base: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     idleConnTimeout,
		},
	}
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// use the base RoundTripper to make the request
	return r.base.RoundTrip(req)
}
