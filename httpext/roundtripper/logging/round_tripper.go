package logging

import (
	"context"
	"log"
	"net/http"
	"time"
)

// LoggingHeaderRoundTripper adds a header and logs request duration.
type LoggingHeaderRoundTripper struct {
	HeaderName  string            // Name of the header to add
	HeaderValue string            // Value of the header to add
	Proxied     http.RoundTripper // The next RoundTripper in the chain
}

// NewLoggingHeaderRoundTripper creates a new LoggingHeaderRoundTripper.
// If base is nil, http.DefaultTransport is used.
func NewLoggingHeaderRoundTripper(headerName, headerValue string, base http.RoundTripper) *LoggingHeaderRoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &LoggingHeaderRoundTripper{
		HeaderName:  headerName,
		HeaderValue: headerValue,
		Proxied:     base,
	}
}

// RoundTrip adds a header, executes the request using the Proxied RoundTripper,
// and logs the duration.
func (rt *LoggingHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original request headers
	clonedReq := req.Clone(req.Context())

	// Add the custom header
	if rt.HeaderName != "" {
		clonedReq.Header.Set(rt.HeaderName, rt.HeaderValue)
	}

	log.Printf("Sending request to %s with header '%s: %s'", clonedReq.URL, rt.HeaderName, rt.HeaderValue)

	start := time.Now()
	// Call the *next* RoundTripper in the chain (e.g., http.DefaultTransport)
	// DO NOT call back into a client's Do method here.
	resp, err := rt.Proxied.RoundTrip(clonedReq)
	duration := time.Since(start)

	if err != nil {
		log.Printf("Request to %s failed after %v: %v", clonedReq.URL, duration, err)
	} else {
		log.Printf("Request to %s completed in %v with status %s", clonedReq.URL, duration, resp.Status)
	}

	return resp, err
}

// --- Example Usage ---
func Executer() {
	// Create the custom RoundTripper, wrapping the default transport
	customTransport := NewLoggingHeaderRoundTripper("X-Custom-ID", "my-request-123", nil) // nil uses http.DefaultTransport

	// Create an http.Client that uses your custom RoundTripper
	myClient := &http.Client{
		Transport: customTransport,
		Timeout:   30 * time.Second, // Overall client timeout
	}

	// Create a request with its own context timeout (optional, but recommended)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Per-request timeout
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/2", nil) // Example URL
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Use the client
	resp, err := myClient.Do(req)
	if err != nil {
		log.Printf("Client.Do error: %v", err)
		// Handle error (e.g., context deadline exceeded, connection error)
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response: %s", resp.Status)
	// Process response...
}
