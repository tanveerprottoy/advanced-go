package squarercancellation

import (
	"fmt"
	"log"
	"sync"
)

// There’s no formal definition of a pipeline in Go;
// it’s just one of many kinds of concurrent programs.
// Informally, a pipeline is a series of stages connected by channels,
// where each stage is a group of goroutines running the same function.
// In each stage, the goroutines:

// receive values from upstream via inbound channels
// perform some function on that data, usually producing new values
// send values downstream via outbound channels

// Each of the pipeline stages is free to return as soon as done is closed.

// Here are the guidelines for pipeline construction:
// stages close their outbound channels when all the send operations are done.
// stages keep receiving values from inbound channels until those channels are closed
// or the senders are unblocked.

// The first stage, gen, is a function that converts a list of integers to
// a channel that emits the integers in the list. The gen function starts
// a goroutine that sends the integers on the channel and closes the channel
// when all the values have been sent. it stops when value is received on done
func gen(done <-chan struct{}, nums ...int) <-chan int {
	// use buffered channel
	out := make(chan int)

	go func() {
		defer close(out)

		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()

	return out
}

// The second stage, sq, receives integers from a channel and returns a channel
// that emits the square of each received integer. After the inbound channel
// is closed and this stage has sent all the values downstream, it closes
// the outbound channel
// sq can return as soon as done is closed. sq ensures its out channel is closed on
// all return paths via a defer statement
func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()

	return out
}

// sink/consumer function that receives values from the second
// stage and prints each one,
// until the channel is closed or done channel is closed
func sink(done <-chan struct{}, in <-chan int) {
	for n := range in {
		select {
		case <-done:
			return
		default:
			log.Println(n)
		}
	}
}

// merge reads from multiple channels and combine them
// this is a fan in pattern
// A function can read from multiple inputs and proceed until all are closed by multiplexing
// the input channels onto a single channel that’s closed when all the inputs are closed.
// This is called fan-in.
// The output routine in merge can return without draining its inbound channel,
// since it knows the upstream sender, sq, will stop attempting to send when done is
// closed. output ensures wg.Done is called on all return paths via a defer statement:
func merge(done chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c or done is closed, then calls
	// wg.Done.
	output := func(c <-chan int) {
		defer wg.Done()

		for n := range c {
			// a select statement that proceeds either when the send on out happens
			// or when they receive a value from done, in case of done we return
			// the func and the deferred wg.done is called
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	// This must start after the wg.Add call.
	go func() {
		wg.Wait()

		close(out)
	}()

	return out
}

func Executer() {
	// We need a way to tell an unknown and unbounded number of goroutines to stop sending
	// their values downstream. In Go, we can do this by closing a channel, because
	// a receive operation on a closed channel can always proceed immediately, yielding
	// the element type’s zero value.
	// This means that main can unblock all the senders simply by closing the done channel.
	// This close is effectively a broadcast signal to the senders. We extend each of our
	// pipeline functions to accept done as a parameter and arrange for the close to happen
	// via a defer statement, so that all return paths from main will signal the pipeline
	// stages to exit.

	// Set up a done channel that's shared by the whole pipeline,
	// and close that channel when this pipeline exits, as a signal
	// for all the goroutines we started to exit.
	done := make(chan struct{})
	defer close(done)

	in := gen(done, 2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(done, in)
	c2 := sq(done, in)

	// Consume the first value from output.
	out := merge(done, c1, c2)
	fmt.Println(<-out) // 4 or 9

	// done will be closed by the deferred call.
}
