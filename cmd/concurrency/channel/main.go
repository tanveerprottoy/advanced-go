package main

import (
	"fmt"

	"github.com/tanveerprottoy/concurrency-go/internal/app/concurrency/channel"
)

func main() {
	/* ch := make(chan int)
	channel.Receive(ch)
	// lock main goroutine
	<-ch
	fmt.Println(ch) */

	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	// split the slice in half
	// spawn a go routine for
	// each half
	go channel.Sum(s[:len(s)/2], c)
	go channel.Sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c
	fmt.Println(x, y, x+y)
}
