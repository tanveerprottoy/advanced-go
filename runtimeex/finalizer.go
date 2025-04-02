package runtimeex

// SetFinalizer sets the finalizer associated with obj to the provided finalizer 
// function. When the garbage collector finds an unreachable block with an associated 
// finalizer, it clears the association and runs finalizer(obj) in a separate goroutine.
// This makes obj reachable again, but now without an associated finalizer. Assuming 
// that SetFinalizer is not called again, the next time the garbage collector sees that 
// obj is unreachable, it will free obj.

// SetFinalizer(obj, nil) clears any finalizer associated with obj.

// New Go code should consider using AddCleanup instead, which is much less error-prone 
// than SetFinalizer.