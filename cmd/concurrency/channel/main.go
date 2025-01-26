package main

import (
	"fmt"
	"log"

	"github.com/tanveerprottoy/advanced-go/concurrency/channel"
)

func sumSplit(vals []int) {
	ch := make(chan int)
	// split the slice in half
	// spawn a goroutine for
	// each half
	go channel.Sum(vals[:len(vals)/2], ch)
	go channel.Sum(vals[len(vals)/2:], ch)
	// receive from ch
	x, y := <-ch, <-ch
	fmt.Printf("split sum: half0 - %d, half1 - %d, result - %d \n", x, y, x+y)
}

func sumMultiply(vals []int) {
	ch := make(chan int)
	ch2 := make(chan int)
	// spawn two goroutines for sum and multiply
	go channel.Sum(vals, ch)
	go channel.Multiply(vals, ch2)
	// receive from ch
	x, y := <-ch, <-ch2
	fmt.Printf("sum & multiply: sum - %d, multiply - %d", x, y)
}

func pingPong() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	go channel.Ping(pings, "passed message")
	go channel.Pong(pings, pongs)
	fmt.Println(<-pongs)
}

func worker() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	go channel.Worker(jobs, results)
	go channel.Worker(jobs, results)
	go channel.Worker(jobs, results)
	go channel.Worker(jobs, results)
	go channel.Worker(jobs, results)

	for i := 0; i < 100; i++ {
		jobs <- i
	}
	close(jobs)
	for j := 0; j < 100; j++ {
		log.Println(<-results)
	}
}

func main() {
	/* ch := make(chan int)
	channel.Receive(ch)
	// block main goroutine
	<-ch
	// will execute below line when ch
	// receives a value
	fmt.Println(ch) */
	vals := []int{7, 2, 8, 9, 4, 5}
	sumSplit(vals)
	sumMultiply(vals)
	pingPong()
	worker()
	// channel.Process(s, bch)
}
