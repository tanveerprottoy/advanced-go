package channel

import "fmt"

func Receive(ch chan int) {
	fmt.Println(ch)
	go func(ch chan int) {
		ch <- 2
	}(ch)
}

func Sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}
