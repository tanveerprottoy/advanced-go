/*
package leak demonstrates some example and patterns of leak that might happen in Go
*/
package leak

import "os"

// Go reclaims memory allocated on the heap via garbage collection,
// but this can be a very expensive process. Compare that to the stack where
// memory is "cheap" and is freed automatically when its stack frame is destroyed.
// In order to allocate memory on the stack the Go compiler must evaluate several,
// determining factors:
// pointers to stack objects cannot be stored in the heap
// pointers to stack objects cannot outlive the object's stack frame
// stack objects cannot exceed the size of the stack, ex. a 15 MiB buffer [15 * 1024 * 1024]byte
// The compile-time process Go use to determine whether memory is dynamically
// managed on the heap or can be allocated on the stack is known as escape analysis.
// Escape analysis walks a program's abstract syntax tree (AST) to build a graph of
// all the variables encountered.
// It is possible to see which variables end up on the heap by using the compiler
// flag -m when building (or testing) Go code.
// Ex: go build -gcflags=-m=3 [package]

// Imagine for a moment there is a kitchen sink with a crack in it.
// The sink has the potential to leak, but nothing will escape the basin until
// the sink is used. Much like our imaginary sink, escape analysis inspects variables
// with the potential to escape when they are used. If that potential exists, the
// variable is marked as "leaking."
// Criteria: what are the criteria for leaking?
// Leak destination: where does the water go when it goes down the drain?
// Leaking (to a sink): bye-bye
// Leaking to result: you'll be back!
// Leak without escape: you never left!

// Criteria
// There are two requirements to be eligible for leaking:
// The variable must be a function parameter
// The variable must be a reference type, ex. channels, interfaces, maps, pointers, slices
// Value types such as built-in numeric types, structs, and arrays are not elgible to be leaked. That does not mean they are never placed on the heap, it just means a parameter of int32 is not going to send you running for a mop anytime soon.
// If the above criteria is met, then a parameter will leak if:
// The variable is returned from the same function and/or
// is assigned to a sink outside of the stack frame to which the variable belongs.

// Leak destination
// There are two primary types of leaks:

// leaking (to a sink)
// leaking to result
// Leaking (to a sink)
// If a function's parameter is a reference type and the function assigns the parameter
// to a variable outside of the function, the variable is leaking. While the compiler
// flag -m that prints optimizations does not indicate to where the parameter is leaking,
// it is helpful to think of this as _leaking to a sink
// Ex:
// The function leakToSink leaks the parameter id to the package-level field sink.
var sink *int32

func leakToSink(id *int32) { // leaking param: id
	sink = id
}

// Leaking to result
// Another type of leak is when a reference parameter is returned from a function:
// Other than the fact that validateID has very poor validation logic (indeed some
// might call it non-existent 😃), the function returns the id parameter. Because
// the value of id is returned, it means it outlives the function's stack frame and
// has the potential to escape to the heap. Therefore the parameter is marked
// as leaking to result.
func leakToResult(id *int32) *int32 { // leaking param: id to result ~r1 level=0
	return id
}

// Leak without escape
// Please remember, a leak is about the potential to escape to the heap. For example,
// a parameter can leak to result without ever escaping to the heap. For example:
// The id parameter for the validateID function is leaking to result because the
// function returns the incoming parameter and thus there is potential for the value
// of id to outlive its stack frame.
// The value id1 is not even mentioned because it is a value type and escape analysis
// only applies to reference types. While a pointer to id1 was passed into the
// validateID function, the Go compiler optimized the pointer to the stack.
// The value id2 is a reference type, but it does not escape. Even though the return
// value of validateID is assigned to validID, its object is on the same stack frame
// as id2, thus the latter does not outlive its stack frame. Therefore id2 does not
// escape to the heap.
func leakWithoutEscape() {
	var id1 int32 = 4096

	if leakToResult(&id1) == nil {
		os.Exit(1)
	}

	var id2 *int32 = new(int32) // new(int32) does not escape

	*id2 = 4096

	validID := leakToResult(id2)

	if validID == nil {
		os.Exit(1)
	}
}
