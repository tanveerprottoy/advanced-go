package main

import (
	"fmt"

	"github.com/tanveerprottoy/concurrency-go/internal/app/function"
)

func main() {
	pos, neg := function.Adder(), function.Adder()
	fmt.Println(pos(3))
	// fmt.Println(pos(3))
	fmt.Println(neg(-3))
}
