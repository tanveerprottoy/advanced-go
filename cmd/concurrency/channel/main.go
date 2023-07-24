package main

import (
	"fmt"

	"github.com/tanveerprottoy/concurrency-go/internal/app/concurrency/channel"
)

func main() {
	ch := make(chan int)
	channel.Receive(ch)
	// lock main goroutine
	<-ch
	fmt.Println(ch)
}
