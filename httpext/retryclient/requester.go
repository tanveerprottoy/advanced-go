package retryclient

import (
	"context"
	"io"
	"net/http"
)

type Requester[R, E any] interface {
	Request(
		ctx context.Context,
		method string,
		url string,
		header http.Header,
		body io.Reader,
	) (*R, *E, error)

	RequestRaw(
		ctx context.Context,
		method string,
		url string,
		header http.Header,
		body io.Reader,
	) (int, []byte, error)
}
