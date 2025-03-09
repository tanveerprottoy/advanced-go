package pipeline

import (
	"log"
	"math/rand"
)

// Create a sequence of numbers, double them, then filter low values,
// and finally print them to console.

// produce generates random numbers
// for the passed count value
func produce(count int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _ = range count {
			// send a random generated number
			out <- rand.Intn(100)
		}
	}()

	return out
}

// doubles doubles the values passed through the
// in channel
func double(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for v := range in {
			// send doubled value to out
			out <- v * 2
		}
	}()

	return out
}

// filter filters the values passed through the
// in channel to greater or equal of min
func filter(in <-chan int, min int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for v := range in {
			// only send out values
			// greater or equal to min
			if v >= min {
				// send out
				out <- v
			}
		}
	}()

	return out
}

func ExecuterDoubleAndFilter() {
	p := produce(20)

	d := double(p)

	f := filter(d, 10)

	for v := range f {
		log.Printf("filtered: %d", v)
	}

	// another approach
	// log.Printf("filtered: %d", filter(double(produce(20)), 10))
}
