package pipeline

import "log"

// returns <-chan string for read only
func producer(strings []string) <-chan string {
	out := make(chan string)

	go func() {
		// need to close the channel when done
		// to stop a deadlock
		// otherwise the receiver (sink) will block
		// forever expecting new values even though
		// the sender has finished sending
		defer close(out)

		for _, s := range strings {
			out <- s
		}
	}()

	return out
}

// a sink, or data sink generally refers to the destination of data flow.
// values <-chan string is a read only channel
func sink(values <-chan string) {
	for value := range values {
		log.Println(value)
	}
}

func ExecuterSimple() {
	src := []string{"foo", "bar", "football"}

	out := producer(src)

	sink(out)
}
