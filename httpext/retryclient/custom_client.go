package retryclient

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"
)

// CustomClient is a custom HTTP client that implements the Client interface
type customClient struct {
	maxRetries int

	httpClient *http.Client
}

func NewCustomClient(maxIdleConnsPerHost, maxRetries int, timeout time.Duration, checkRedirectFunc func(req *http.Request, via []*http.Request) error) *customClient {
	c := &customClient{
		maxRetries: maxRetries,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
			},
		},
	}

	if checkRedirectFunc != nil {
		c.httpClient.CheckRedirect = checkRedirectFunc
	}

	return c
}

func (c *customClient) backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

func (c *customClient) jitter(max, min, attempts int) int {

}

func (c *customClient) restoreRequestBody(req *http.Request) error {
	if req.Body == nil {
		return nil
	}

	// Read the body into a buffer
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// Restore the body so it can be read again
	// This is important because the body is an io.ReadCloser and can only be read once
	req.Body = io.NopCloser(bytes.NewBuffer(buf))

	return nil
}

func (c *customClient) shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return true
	}
	return false
}

func (c *customClient) isRetryableError(err error) bool {
	if err != nil {
		return false
	}

	// check if error is temporary
	if errNet, ok := err.(interface{ Temporary() bool }); ok && errNet.Temporary() {
		return true
	}

	return false
}

func (c *customClient) drainBody(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (c *customClient) Request(req *http.Request) (*http.Response, error) {
	var attempts int

	var (
		res         *http.Response
		err         error
		shouldRetry bool
	)

	for attempts < c.maxRetries {
		// increment attempts
		attempts++
		
		// reusing a request body can be a bit tricky because the
		// io.ReadCloser interface, which is the type of r.Body in an
		// http.Request, is designed for single consumption. Once you've read the body, the underlying reader is often at its end, and attempting to read it again will yield an empty result or an error.
		// Always rewind/restore the request body when non-nil.
		if req.Body != nil {
			if err := c.restoreRequestBody(req); err != nil {
				return nil, err
			}
		}

		res, err = c.httpClient.Do(req)
		if err != nil {
			// check if error is temporary
			if c.isRetryableError(err) {
				// wait for backoff time
				time.Sleep(c.backoff(attempts + 1))
			}

			return nil, err
		}
	}

	defer c.httpClient.CloseIdleConnections()

	return res, nil
}

/* func (c *CustomClient) RequestWithContextRetry(ctx context.Context, method string, url string, header http.Header, body io.Reader, retryCount int) (*http.Response, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.HTTPClient.Timeout)
		defer cancel()
	}

	var (
		req *http.Request
		res *http.Response
		err error
	)

	for i := range retryCount {
		req, err = http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			// failed to create request
			// stop and return error
			return nil, err
		}

		if header != nil {
			req.Header = header
		}

		res, err = c.HTTPClient.Do(req)
		if err != nil {
			if c.shouldRetry(res, err) {
				// wait for backoff time
				time.Sleep(c.backoff(i + 1))
				continue
			}

			return nil, err
		}
	}

	return res, nil
} */

// ex:
// resp, err := http.PostForm("http://example.com/form",
// url.Values{"key": {"Value"}, "id": {"123"}})
func (c *customClient) PostForm(url string, header http.Header, values url.Values) (*http.Response, error) {
	res, err := http.PostForm(url, values)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *customClient) HTTPClient() *http.Client {
	return c.httpClient
}

// BodyBytes allows accessing the request body. It is an analogue to
// http.Request's Body variable, but it returns a copy of the underlying data
// rather than consuming it.

// This function is not thread-safe; do not call it at the same time as another
// call, or at the same time this request is being used with Client.Do.
func (c *customClient) BodyBytes(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	// Read the body into a buffer
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Restore the body so it can be read again
	r.Body = io.NopCloser(bytes.NewBuffer(buf))

	return buf, nil
}
