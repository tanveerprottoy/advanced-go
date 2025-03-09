package pipeline

import (
	"context"
	"log"
	"time"
)

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
			// simulate long running operation
			time.Sleep(1 * time.Second)

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
	src := []string{"foo", "bar", "football", "tennis", "motogp", "f1"}

	out := producer(src)

	consumer(out)
}

func producerCtx(ctx context.Context, strings []string) <-chan string {
	out := make(chan string)

	go func() {
		// need to close the channel when done
		// to stop a deadlock
		// otherwise the receiver (sink) will block
		// forever expecting new values even though
		// the sender has finished sending
		defer close(out)

		for _, s := range strings {
			// simulate long running operation
			time.Sleep(1 * time.Second)

			// need to check on context
			// to check for cancellation
			// A select blocks until one of its cases can run, then it executes that case.
			// It chooses one at random if multiple are ready.
			// here this select will not block as it try the first case if it's not run
			// then it will go to the next case as it will always run, sending a value to out
			select {
			case <-ctx.Done():
				return
			case out <- s:
			}
		}
	}()

	return out
}

// a consumer, or data consumer generally refers to the destination of data flow
// it's the last stage of a pipeline
// values <-chan string is a read only channel
func consumerCtx(ctx context.Context, values <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-values:
			if !ok {
				return
			}

			log.Println(v)
		}
	}
}

func ExecuterSimpleCtx() {
	src := []string{"foo", "bar", "football", "tennis", "motogp", "f1"}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	out := producerCtx(ctx, src)

	consumerCtx(ctx, out)
}
