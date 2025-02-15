package concurrency

import (
	"context"
	"log"
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	ctx := r.Context()
	processTime := time.Duration(11) * time.Second
	select {
	case <-ctx.Done():
		log.Println("request cancelled")
		return
	case <-time.After(processTime):
		log.Println("request processed")
	}

	// prepare the context with timeout
	// Create a Context with a timeout.
	ctx, cancelFn := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancelFn()
	select {
	case <-ctx.Done():
		log.Println("request cancelled")
		return
	default:
		// continue processing the request
	}

}
