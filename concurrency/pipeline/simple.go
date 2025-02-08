package pipeline

import "log"

// a producer aka source is the first stage
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

// a consumer, or data consumer generally refers to the destination of data flow
// it's the last stage of a pipeline
// values <-chan string is a read only channel
func consumer(values <-chan string) {
	for v := range values {
		log.Println(v)
	}
}

func ExecuterSimple() {
	src := []string{"foo", "bar", "football"}

	out := producer(src)

	consumer(out)
}
