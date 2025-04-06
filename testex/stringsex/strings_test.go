package stringsex_test

import (
	"fmt"
	"testing"

	"github.com/tanveerprottoy/advanced-go/testex/stringsex"
)

// any benchmark should be careful to avoid compiler optimisations eliminating the
// function under test and artificially lowering the run time of the benchmark.
var result string

func BenchmarkConcatenateBuffer(b *testing.B) {
	var s string

	for range b.N {
		s = stringsex.ConcatenateBuffer("Hello ", "World")
	}

	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = s
}

func BenchmarkConcatenateJoin(b *testing.B) {
	var s string

	for range b.N {
		s = stringsex.ConcatenateJoin("Hello ", "World")
	}

	// always store the result to a package level variable
	// so the compiler cannot eliminate the Benchmark itself.
	result = s
}

func BenchmarkConcatenation(b *testing.B) {
	var s string
	lengths := []int{2, 16, 128, 1024, 8192, 65536, 524288, 4194304, 16777216, 134217728}
	for _, l := range lengths {
		first := stringsex.GenerateRandomString(l)
		second := stringsex.GenerateRandomString(l)

		b.Run(fmt.Sprintf("ConcatenateJoin-%d", l), func(b *testing.B) {
			for range b.N {
				s = stringsex.ConcatenateJoin(first, second)
			}

			result = s
		})

		b.Run(fmt.Sprintf("ConcatenateBuffer-%d", l), func(b *testing.B) {
			for range b.N {
				s = stringsex.ConcatenateBuffer(first, second)
			}

			result = s
		})
	}
}
