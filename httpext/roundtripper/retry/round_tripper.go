package retry

import (
	"bytes"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// RoundTripper is a custom HTTP round tripper that implements the http.RoundTripper interface
// Roundtripper should be used when you want to add the retry logic in the http client's
// Transport/Roundtripper level, instead of the client level
type RoundTripper struct {
	maxRetries int
	maxJitter  int

	base http.RoundTripper
}

func NewRoundTripper(maxRetries, maxIdleConnsPerHost int, idleConnTimeout time.Duration) *RoundTripper {
	return &RoundTripper{
		maxRetries: maxRetries,
		base: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     idleConnTimeout,
		},
	}
}

// backoff generates the backoff time in seconds based on the number of retries
func (r *RoundTripper) backoff(retries int) time.Duration {
	// 2^n backoff, n = number of retries
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

// jitter generates a random jitter in milliseconds which is added to the backoff time
func (r *RoundTripper) jitter(max, attempts int) time.Duration {
	rnd := rand.Intn(max)

	return time.Duration(int(attempts*rnd)) * time.Millisecond
}

func (r *RoundTripper) restoreRequestBody(req *http.Request) error {
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

func (r *RoundTripper) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// check if error is temporary
	if errNet, ok := err.(interface{ Temporary() bool }); ok && errNet.Temporary() {
		return true
	}

	return false
}

func (r *RoundTripper) drainBody(resp *http.Response) {
	// drain the response body to reuse the connection
	// only do this if the response is not nil and the body is not nil
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println("RoundTripper: RoundTrip called")

	var (
		attempts int
		resp     *http.Response
		err      error
	)

	for attempts < r.maxRetries {
		log.Printf("customClient: doWithRetry called, attempt %d\n", attempts)

		// reusing a request body can be a bit tricky because the
		// io.ReadCloser interface, which is the type of r.Body in an
		// http.Request, is designed for single consumption. Once you've read the body, the underlying reader is often at its end, and attempting to read it again will yield an empty result or an error.
		// Always rewind/restore the request body when non-nil.
		if req.Body != nil {
			if err := r.restoreRequestBody(req); err != nil {
				return nil, err
			}
		}

		// use the base RoundTripper to make the request
		resp, err = r.base.RoundTrip(req)

		log.Printf("customClient: doWithRetry.httpClient.Do called, attempt %d, response: %v\nerr: %v\n", attempts, resp, err)

		if err != nil {
			log.Printf("customClient: doWithRetry.httpClient.Do if err != nil: %v\n", err)

			// check if error is temporary
			if r.isRetryableError(err) {
				// temp
				log.Printf("backoff: %v\njitter: %v\n", r.backoff(attempts), r.jitter(r.maxJitter, attempts))

				// drain the response body to reuse the connection
				r.drainBody(resp)

				// wait for backoff time
				time.Sleep(r.backoff(attempts) + r.jitter(r.maxJitter, attempts))

				// continue to next attempt
				continue
			}

			return nil, err
		}

		// increment attempts
		attempts++

		log.Printf("customClient: doWithRetry after increment attempts: %d\n", attempts)
	}

	return resp, nil
}
