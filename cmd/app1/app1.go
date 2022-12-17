package main

import "fmt"

// limit goroutine with buffer channel
func channelBuffer() {
	ch := make(chan int, 15)
	/*ch <- v    // Send v to channel ch.
	  v := <-ch  // Receive from ch, and
           		 // assign value to v. 
	*/
	ch <- 1
	ch <- 2
	v := <-ch
	fmt.Println(v)
	fmt.Println(<-ch)
}

func main() {
	go channelBuffer()
}
