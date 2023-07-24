package channel

import "fmt"

func Adder() func(int) int {
	sum := 2
	return func(x int) int {
		// fmt.Println(x)
		sum += x
		return sum
	}
}

func ReturnAFunc(str string) func(string) string {
	str += "!"
	return func(str string) string {
		return str
	}
}

func Receive(ch chan int) {
	fmt.Println(ch)
	go func(ch chan int) {
		ch <- 2
	}(ch)
}
