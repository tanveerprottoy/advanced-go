package function

import "fmt"

// closure example
func Adder() func(int) int {
	sum := 2
	return func(x int) int {
		fmt.Println("func param x: ", x)
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
