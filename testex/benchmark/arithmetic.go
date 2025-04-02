package benchmark

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Add[T Numeric](a, b T) T {
	return a + b
}

func Subract[T Numeric](a, b T) T {
	return a - b
}

// The factorial of n is denoted by n! and calculated by multiplying the integer 
// numbers from 1 to n. The formula for n factorial is n! = n × (n - 1)!.
func Factorial(n int) int {
	if n == 0 {
		return 1
	}
	
	return n * Factorial(n-1)
}