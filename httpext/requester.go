package httpext

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
		retry bool,
	) (*R, *E, error)

	RequestRaw(
		ctx context.Context,
		method string,
		url string,
		header http.Header,
		body io.Reader,
		retry bool,
	) (int, []byte, error)
}
