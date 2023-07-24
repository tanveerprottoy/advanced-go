package main

import (
	"fmt"

	"github.com/tanveerprottoy/concurrency-go/internal/channel"
)

func main() {
	/* ch := make(chan int)
	channel.Receive(ch)
	// lock main goroutine
	<-ch
	fmt.Println(ch) */

	pos, neg := channel.Adder(), channel.Adder()
	fmt.Println(pos(3))
	// fmt.Println(pos(3))
	fmt.Println(neg(-3))
}
