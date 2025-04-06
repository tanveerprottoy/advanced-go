package httpext

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// retryHTTP performs an HTTP request with retries, using exponential backoff and jitter.
func retryHTTP(ctx context.Context, client *http.Client, req *http.Request, maxRetries int, retryDelay time.Duration) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = client.Do(req.WithContext(ctx))
		if err == nil {
			// Request successful.
			return resp, nil
		}

		// Check for context cancellation.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Check for retryable errors.  Adjust as needed for your application.
		if isRetryableError(err) {
			// Apply exponential backoff with jitter.
			delay := retryDelay * time.Duration(attempt)
			jitter := time.Duration(rand.Int63n(int64(delay)))
			time.Sleep(delay + jitter)
			log.Printf("Retrying request (attempt %d/%d): %v", attempt, maxRetries, err)
			continue
		}

		// Non-retryable error.
		return nil, fmt.Errorf("request failed after %d retries: %w", attempt-1, err)
	}

	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

// isRetryableError checks if the given error is retryable.  This is a sample; customize as needed.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	// Check for common transient network errors.
	if netErr, ok := err.(interface{ Temporary() bool }); ok && netErr.Temporary() {
		return true
	}
	// Add other retryable error conditions as needed.  For example, specific HTTP status codes.
	return false
}

func requestWithRetry() {
	// Example usage:
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := retryHTTP(context.Background(), client, req, 3, 500*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Request successful: Status Code %d\n", resp.StatusCode)
}

// retryTransport wraps an http.RoundTripper to add retry logic.
type retryTransport struct {
	baseTransport http.RoundTripper
	retries       int
	delay         time.Duration
}

// RoundTrip implements the retry logic for HTTP requests.
func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rt.retries; attempt++ {
		// Perform the HTTP request
		resp, err = rt.baseTransport.RoundTrip(req)
		if err == nil && resp.StatusCode < 500 {
			// Success or non-retriable status code
			return resp, nil
		}

		// Log the retry attempt
		log.Printf("Retry attempt %d for %s (error: %v, status: %v)\n", attempt+1, req.URL, err, resp.StatusCode)

		// Close the response body to avoid resource leaks
		if resp != nil {
			resp.Body.Close()
		}

		// Wait before retrying
		time.Sleep(rt.delay)
	}

	// Return the last error if all retries fail
	if err == nil {
		err = errors.New("max retries exceeded")
	}
	return nil, err
}

// newRetryClient creates an HTTP client with retry logic.
func newRetryClient(retries int, delay time.Duration) *http.Client {
	return &http.Client{
		Transport: &retryTransport{
			baseTransport: http.DefaultTransport,
			retries:       retries,
			delay:         delay,
		},
	}
}

func requestWithRetry2() {
	client := newRetryClient(3, 2*time.Second) // 3 retries with 2 seconds delay

	req, err := http.NewRequest("GET", "https://httpstat.us/500", nil) // Simulate a server error
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Request succeeded with status: %d\n", resp.StatusCode)
}
