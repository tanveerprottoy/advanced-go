package channel

import "fmt"



func Receive(ch chan int) {
	fmt.Println(ch)
	go func(ch chan int) {
		ch <- 2
	}(ch)
}
