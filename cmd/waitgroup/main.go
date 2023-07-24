package main

import (
	"fmt"

	"github.com/tanveerprottoy/code-practise/concurrency"
)

func main() {
	ch := make(chan int)
	go concurrency.Subtract(5, 3, ch)
	o := <-ch
	fmt.Println(o)
	go concurrency.Multiply(o, 2, ch)
	o = <-ch
	fmt.Println(o)
}
