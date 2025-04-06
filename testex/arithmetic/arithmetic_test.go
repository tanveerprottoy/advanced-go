package arithmetic_test

import (
	"math/rand"
	"testing"

	"github.com/tanveerprottoy/advanced-go/testex/arithmetic"
)

// tests
func TestAdd(t *testing.T) {
	// run parallelly
	t.Parallel()

	tests := []struct {
		name string
		val0 int
		val1 int
		exp  int
	}{
		{"2 + 5", 2, 5, 7},
		{"9 + 5", 9, 5, 14},
		{"27 + 45", 27, 45, 72},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := arithmetic.Add(tc.val0, tc.val1)

			if actual != tc.exp {
				t.Errorf("Add(%d, %d) = %v; want %v", tc.val0, tc.val1, actual, tc.exp)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	// run parallelly
	t.Parallel()

	tests := []struct {
		name string
		val0 int
		val1 int
		exp  int
	}{
		{"5 - 3", 5, 3, 2},
		{"9 - 5", 9, 5, 4},
		{"7 - 5", 7, 5, 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := arithmetic.Subract(tc.val0, tc.val1)

			if actual != tc.exp {
				t.Errorf("Subract(%d, %d) = %v; want %v", tc.val0, tc.val1, actual, tc.exp)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	// run parallelly
	t.Parallel()

	tests := []struct {
		name string
		in   int
		exp  int
	}{
		{"2!", 2, 2},
		{"3!", 3, 6},
		{"4!", 4, 24},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := arithmetic.Factorial(tc.in)

			if actual != tc.exp {
				t.Errorf("Factorial(%d) = %v; want %v", tc.in, actual, tc.exp)
			}
		})
	}
}

// bencmarks
// any benchmark should be careful to avoid compiler optimisations eliminating the
// function under test and artificially lowering the run time of the benchmark.
const max = 100

var result int

func BenchmarkAdd(b *testing.B) {
	var r int

	// good practise not to use values from b.N
	for range b.N {
		r = arithmetic.Add(rand.Intn(max), rand.Intn(max))
	}

	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = r
}

func BenchmarkSubtract(b *testing.B) {
	var r int

	// good practise not to use values from b.N
	for range b.N {
		r = arithmetic.Subract(rand.Intn(max), rand.Intn(max))
	}

	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = r
}

func BenchmarkFactorial(b *testing.B) {
	var r int

	// good practise not to use values from b.N
	for range b.N {
		r = arithmetic.Factorial(rand.Intn(max))
	}

	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = r
}

func benchmarkFactorial(i int, b *testing.B) {
	var r int

	// run the factorial function
	for range b.N {
		r = arithmetic.Factorial(i)
	}

	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = r
}

func BenchmarkFactorial5(b *testing.B) {
	benchmarkFactorial(5, b)
}

func BenchmarkFactorial10(b *testing.B) {
	benchmarkFactorial(10, b)
}

func BenchmarkFactorial15(b *testing.B) {
	benchmarkFactorial(15, b)
}

func BenchmarkFactorial24(b *testing.B) {
	benchmarkFactorial(24, b)
}

func BenchmarkAll(b *testing.B) {
	b.Run("Add", BenchmarkAdd)

	b.Run("Subtract", BenchmarkSubtract)

	b.Run("Factorial", BenchmarkFactorial)
}
