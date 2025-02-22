package contextex

import (
	"context"
	"log"
	"math/rand"
	"time"
)

// establishNetworkConnection simulates a long running network connection
// a write only channel is passed to the function
func establishNetworkConnection(ch chan<- bool) {
	log.Println("Establishing network connection...")

	// simulate success and failure
	// based on random number
	if rand.Intn(10) < 5 {
		time.Sleep(2 * time.Second)

		log.Println("Network connection established")
	} else {
		log.Println("Failed to establish network connection")

		// simulate hanging long time
		time.Sleep(10 * time.Minute)
	}

	// write/send value to the channel
	ch <- true
}

func establishNetworkConnectionWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	defer cancel()

	// create a channel to receive the result
	// it's an unbuffered channel
	ch := make(chan bool)

	// start the goroutine
	go establishNetworkConnection(ch)

	// wait on multiple channels
	select {
	case <-ch:
		log.Println("Network connection established")
	case <-ctx.Done():
		log.Println("Timeout. Cancelling network connection")
	}

}

func ExecuterEstablishNetConnection() {
	establishNetworkConnectionWithTimeout()
}
