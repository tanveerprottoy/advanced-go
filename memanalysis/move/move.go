package move

// Move
// We have discussed when parameters are leaked and when variables escape to the heap,
// but what about when the heap is so lovely a variable decides to move there? This
// page discusses when and why escape analysis moves some variables to the heap.

// Criteria: what are the criteria for moving to the heap?
// Moving to the heap: when a value on the stack is moved to the heap
// Storing value types in interfaces: when a value type escapes to the heap after all
// Criteria
// There are two requirements to be eligible for moving to the heap:

// The variable must be a value type
// A reference to the variable is assigned to a location outside of the local stack frame
// What does this look like? It is actually semantically the same as when a variable
// escapes to the heap -- it all just depends on the initial declaration
// of that varaible's type...

// Moving to the heap
// A value type variable is moved to the heap if a reference to it is assigned to a
// location outside of the variable's local stack frame.

// Unlike when id := new(int32) escapes to the heap, the id in the example is moved
// to the heap because it was initially a value type. Thus, if all other conditions
// are met for escaping and moving to the heap:

// a value escapes if it is initially a reference type
// a value moves if it is initially a value type
// Actually, there is one exception to this rule...
var sink *int32

//go:noinline
func moveToHeap() {
	var id int32 = 4096
	sink = &id
}

// Actually, there is one exception to this rule...
// continued...
// Storing value types in interfaces
// It would be too simple if reference types escaped to the heap and value types moved,
// wouldn't it? One way a value type is marked as escaping to the heap is when the
// value is stored in an interface. For example:
var sink2 interface{}

//go:noinline
func moveToHeap2() {
	var id int32 = 4096
	sink2 = id
}
