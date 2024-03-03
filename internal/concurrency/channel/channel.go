package channel

import (
	"fmt"
)

func Receive(ch chan int) {
	fmt.Println(ch)
	go func(ch chan int) {
		ch <- 2
	}(ch)
}

func Sum(vals []int, ch chan int) {
	res := 0
	if len(vals) == 0 {
		ch <- res
		return
	}
	for _, v := range vals {
		res += v
	}
	// send res to ch
	ch <- res
}

func Multiply(vals []int, ch chan int) {
	res := 1
	if len(vals) == 0 {
		ch <- res
		return
	}
	for _, v := range vals {
		res *= v
	}
	// send res to ch
	ch <- res
}

// When using channels as function parameters,
// you can specify if a channel is meant to only send or receive values.
// This specificity increases the type-safety of the program

// This ping function only accepts a channel for sending values.
// It would be a compile-time error to try to receive on this channel
func Ping(pings chan<- string, msg string) {
	pings <- msg
}

// The pong function accepts one channel for receives (pings)
// and a second for sends (pongs).
func Pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-1)
}

func Fibonacci2(n int, ch chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		ch <- x
		x, y = y, x+y
	}
	close(ch)
}

// channel direction
func Worker(jobs <-chan int, results chan<- int) {
	for j := range jobs {
		results <- Fibonacci(j)
	}
}

func Process(vals []int, ch chan int) {
	ch <- 5
}
