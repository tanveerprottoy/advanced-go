package httpext

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// service implements the Requester interface
// it makes http requests using the client
type service[R, E any] struct {
	client Client
}

func NewService[R, E any](client Client) *service[R, E] {
	return &service[R, E]{
		client: client,
	}
}

func (s *service[R, E]) buildRequest(
	ctx context.Context,
	method, url string,
	header http.Header,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

// Request is a generic method to make a request with context
// generic paramters are provided by the struct itself
// Generic parameters: R = response type, E = error type
// use this function when you want to parse the response body to a specific type
// and also parse the error response to a specific type
func (s *service[R, E]) Request(
	ctx context.Context,
	method string,
	url string,
	header http.Header,
	body io.Reader,
	retry bool,
) (*R, *E, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.client.HTTPClient().Timeout)
		defer cancel()
	}

	req, err := s.buildRequest(ctx, method, url, header, body)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, retry)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		// resp ok, parse response body to type
		var r R

		err := json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			return nil, nil, err
		}

		return &r, nil, nil
	} else {
		// resp not ok, parse error
		var e E

		err := json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, nil, err
		}

		return nil, &e, errors.New("error response was returned")
	}
}

// RequestRaw is a function to make a request with context
// it returns the status code and the response body as a byte slice
// use this function when need to get the raw response body
func (s *service[R, E]) RequestRaw(
	ctx context.Context,
	method string,
	url string,
	header http.Header,
	body io.Reader,
	retry bool,
) (int, []byte, error) {
	req, err := s.buildRequest(ctx, method, url, header, body)
	if err != nil {
		return 0, nil, err
	}

	resp, err := s.client.Do(req, retry)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, b, nil
}
