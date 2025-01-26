package goroutine

import (
	"fmt"
	"time"
)

func SleeperFunc(ch chan int, d time.Duration) {
	fmt.Println("sleeping for: ", d)
	time.Sleep(d)
	// send back the channel with value
	ch <- 0
}

func CallerFunc() {
	ch := make(chan int)

	// want to wait and get some value
	// from the SleeperFunc goroutine
	// for that sending the channel
	go SleeperFunc(ch, time.Second*5)

	// receive the value from channel
	res := <-ch
	
	fmt.Println(res)
}
