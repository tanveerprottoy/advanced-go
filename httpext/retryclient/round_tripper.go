package retryclient

import (
	"net/http"
	"sync"
)

// RoundTripper is a custom HTTP round tripper that implements the http.RoundTripper interface
type RoundTripper struct {
	client Client

	once sync.Once
}

func NewRoundTripper() *RoundTripper {
	return new(RoundTripper)
}

func (r *RoundTripper) Init(client Client) {
	r.once.Do(func() {
		r.client = client
	})
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, nil
}
