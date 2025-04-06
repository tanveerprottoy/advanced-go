package arithmetic

import (
	"log"
	"testing"
)

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

var result int

func RunBenchmarkTests() {
	res := testing.Benchmark(func(b *testing.B) {
		var r int

		for i := range b.N {
			r = Add(i, i+1)
		}

		// always store the result to a package level variable
		// so the compiler cannot eliminate the Benchmark itself.
		result = r
	})

	log.Printf("Memory allocations : %d \n", res.MemAllocs)
	log.Printf("Number of bytes allocated: %d \n", res.Bytes)
	log.Printf("Number of run: %d \n", res.N)
	log.Printf("Time taken: %s \n", res.T)
}
