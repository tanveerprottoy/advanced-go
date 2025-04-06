package httptraceex

import (
	"context"
	"log"
	"net/http"
	"net/http/httptrace"
)

// HTTP events
// The httptrace package provides a number of hooks to gather information during an
// HTTP round trip about a variety of events. These events include:
// Connection creation
// Connection reuse
// DNS lookups
// Writing the request to the wire
// Reading the response

// Tracing events
// You can enable HTTP tracing by putting an *httptrace.ClientTrace containing hook
// functions into a request’s context.Context. Various http.RoundTripper
// implementations report the internal events by looking for context’s
// *httptrace.ClientTrace and calling the relevant hook functions.

// The tracing is scoped to the request’s context and users should put a
// *httptrace.ClientTrace to the request context before they start a request.
func trace(ctx context.Context) {
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)

	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Printf("Got Conn: %+v\n", connInfo)
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}

	// During a round trip, http.DefaultTransport will invoke each hook as an event
	// happens. The program above will print the DNS information as soon as the DNS
	// lookup is complete. It will similarly print connection information when a
	// connection is established to the request’s host.
}

// Tracing with http.Client
// The tracing mechanism is designed to trace the events in the lifecycle of a single
// http.Transport.RoundTrip. However, a client may make multiple round trips to
// complete an HTTP request. For example, in the case of a URL redirection, the
// registered hooks will be called as many times as the client follows HTTP redirects,
// making multiple requests. Users are responsible for recognizing such events at the
// http.Client level. The program below identifies the current request by using an
// http.RoundTripper wrapper.

// transport is an http.RoundTripper that keeps track of the in-flight
// request and implements hooks to report HTTP tracing events.
type transport struct {
	current *http.Request
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	return http.DefaultTransport.RoundTrip(req)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *transport) GotConn(info httptrace.GotConnInfo) {
	log.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func traceClient() {
	t := &transport{}

	req, _ := http.NewRequest("GET", "https://google.com", nil)
	trace := &httptrace.ClientTrace{
		GotConn: t.GotConn,
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{Transport: t}

	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
}

func Executer() {
	trace(context.Background())

	traceClient()
}

// The program will follow the redirect of google.com to www.google.com and will output:
// Connection reused for https://google.com? false
// Connection reused for https://www.google.com/? false

// The Transport in the net/http package supports tracing of both HTTP/1 and HTTP/2 requests.
// If you are an author of a custom http.RoundTripper implementation, you can support
// tracing by checking the request context for an *httptest.ClientTrace and invoking
// the relevant hooks as the events occur.
