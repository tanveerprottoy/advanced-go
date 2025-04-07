package retryclient

import (
	"net/http"
	"net/url"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)

	PostForm(url string, header http.Header, values url.Values) (*http.Response, error)

	HTTPClient() *http.Client
}
