package squarerbufferchan

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

// The first stage, gen, is a function that converts a list of integers to
// a channel that emits the integers in the list. The gen function starts
// a goroutine that sends the integers on the channel and closes the channel
// when all the values have been sent:
func gen(nums ...int) <-chan int {
	// use buffered channel
	out := make(chan int, len(nums))

	// when buffered channel is used
	// goroutine is not required to run the
	// loop and close block
	for _, n := range nums {
		out <- n
	}

	close(out)

	return out
}

// The second stage, sq, receives integers from a channel and returns a channel
// that emits the square of each received integer. After the inbound channel
// is closed and this stage has sent all the values downstream, it closes
// the outbound channel:
func sq(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * n
		}

		close(out)
	}()

	return out
}

// sink/consumer function that receives values from the second
// stage and prints each one,
// until the channel is closed
func sink(in <-chan int) {
	for n := range in {
		log.Println(n)
	}
}

// The Executer function sets up the pipeline and runs the final stage:
// it receives values from the second stage and prints each one,
// until the channel is closed:
func ExecuterSquarer() {
	// Set up the pipeline.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	fmt.Println(<-out) // 4
	fmt.Println(<-out) // 9

	// better way to consume the output
	sink(out)
}

func ExLoop() {
	// Set up the pipeline and consume the output.
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n) // 16 then 81
	}
}

// merge reads from multiple channels and combine them
// this is a fan in pattern
// A function can read from multiple inputs and proceed until all are closed by multiplexing
// the input channels onto a single channel that’s closed when all the inputs are closed.
// This is called fan-in.
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int, 1) // enough space for the unread inputs

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}

		wg.Done()
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

func ExecuterSquarerFan() {
	// buffered channel will not block until
	// 1. if the buffered channel is full then send will block
	// 2. if the buffered channel is empty then receive will block
	// c := make(chan int, 2) // buffer size 2
	// c <- 1                 // succeeds immediately
	// c <- 2                 // succeeds immediately
	// c <- 3                 // blocks until another goroutine does <-c and receives 1

	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	// here a single channel is being read by same function multiple times
	// Multiple functions can read from the same channel until that channel is closed;
	// this is called fan-out. This provides a way to distribute work amongst a group of
	// workers to parallelize CPU use and I/O.

	// This approach has a problem: each downstream receiver needs to know the number
	// of potentially blocked upstream senders and arrange to signal those senders on
	// early return. Keeping track of these counts is tedious and error-prone.
	c1 := sq(in)
	c2 := sq(in)

	out := merge(c1, c2)

	fmt.Println(<-out) // 4 or 9

	return

	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}
