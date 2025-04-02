package benchmark_test

import (
	"math/rand"
	"testing"

	"github.com/tanveerprottoy/advanced-go/testex/benchmark"
)

// tests
func TestAdd(t *testing.T) {
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
			actual := benchmark.Add(tc.val0, tc.val1)

			if actual != tc.exp {
				t.Errorf("Add(%d, %d) = %v; want %v", tc.val0, tc.val1, actual, tc.exp)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
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
			actual := benchmark.Subract(tc.val0, tc.val1)

			if actual != tc.exp {
				t.Errorf("Subract(%d, %d) = %v; want %v", tc.val0, tc.val1, actual, tc.exp)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
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
			actual := benchmark.Factorial(tc.in)

			if actual != tc.exp {
				t.Errorf("Factorial(%d) = %v; want %v", tc.in, actual, tc.exp)
			}
		})
	}
}

// benchmarks
func BenchmarkAdd(b *testing.B) {
	for i := range b.N {
		_ = benchmark.Add(rand.Intn(i), rand.Intn(i+2))
	}
}

func BenchmarkSubtract(b *testing.B) {
	for i := range b.N {
		_ = benchmark.Subract(rand.Intn(i), rand.Intn(i+2))
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := range b.N {
		_ = benchmark.Factorial(rand.Intn(i))
	}
}
