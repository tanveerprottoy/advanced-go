package weakpointer

import (
	"runtime"
	"weak"
)

// Just like regular pointers, Weak Pointer may reference any part of an object,
// such as a field of a struct or an element of an array. Objects that are only
// pointed to by weak pointers are not considered reachable, and once the object
// becomes unreachable, Pointer.Value may return nil.
// The primary use-cases for weak pointers are for implementing caches,
// canonicalization maps (like the unique package), and for tying together the
// lifetimes of separate values (for example, through a map with weak keys).
// Two Pointer values always compare equal if the pointers from which they were
// created compare equal. This property is retained even after the object referenced
// by the pointer used to create a weak reference is reclaimed. If multiple weak
// pointers are made to different offsets within the same object (for example,
// pointers to different fields of the same struct), those pointers will not compare
// equal. If a weak pointer is created from an object that becomes unreachable, but is
// then resurrected due to a finalizer, that weak pointer will not compare equal with
// weak pointers created after the resurrection.
// Calling Make with a nil pointer returns a weak pointer whose Pointer.Value always
// returns nil. The zero value of a Pointer behaves as if it were created by passing
// nil to Make and compares equal with such pointers.

// Pointer.Value is not guaranteed to eventually return nil. Pointer.Value may return
// nil as soon as the object becomes unreachable. Values stored in global variables,
// or that can be found by tracing pointers from a global variable, are reachable.
// A function argument or receiver may become unreachable at the last point where
// the function mentions it. To ensure Pointer.Value does not return nil, pass a
// pointer to the object to the runtime.KeepAlive function after the last point
// where the object must remain reachable.
// Note that because Pointer.Value is not guaranteed to eventually return nil, even
// after an object is no longer referenced, the runtime is allowed to perform a
// space-saving optimization that batches objects together in a single allocation
// slot. The weak pointer for an unreferenced object in such an allocation may never
// become nil if it always exists in the same batch as a referenced object. Typically,
// this batching only happens for tiny (on the order of 16 bytes or less) and
// pointer-free objects.

// Weak pointers shine in cases where memory efficiency is crucial. For example:
// Caches: Avoid retaining unused objects longer than necessary.
// Observers: Track objects without preventing their cleanup.
// References: Reduce the risk of memory leaks in long-running programs.
type T struct {
	a int
	b int
}

func ExecuterWP() {
	a := new(string)
	println("original:", a)

	// make a weak pointer
	weakA := weak.Make(a)

	runtime.GC()

	// use weakA
	strongA := weakA.Strong()
	println("strong:", strongA, a)

	runtime.GC()

	// use weakA again
	strongA = weakA.Strong()
	println("strong:", strongA)

	// make a weak pointers
	weakA = weak.Make(a)
	weakA2 := weak.Make(a)

	println("Before GC: Equality check:", weakA == weakA2)

	runtime.GC()

	// Test their equality
	println("After GC: Strong:", weakA.Strong(), weakA2.Strong())
	println("After GC: Equality check:", weakA == weakA2)
}
