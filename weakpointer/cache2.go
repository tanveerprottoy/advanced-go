package weakpointer

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"weak"
)

// Why This Example is Better with Weak Pointers
// Without weak pointers, the cache would hold strong references to all its objects,
// preventing them from being garbage collected. This could lead to memory leaks,
// especially in a long-running server where cached objects accumulate over time.
// By using weak pointers:

// Memory Efficiency: Unused objects are reclaimed by the garbage collector,
// reducing memory usage.
// Automatic Cleanup: You don’t need to implement complex eviction logic.
// Thread Safety: Weak pointers integrate seamlessly into thread-safe structures
// like the Cache in the example.
// Without weak pointers, you’d need a more manual approach, such as periodically
// checking and removing unused objects, which adds complexity and room for bugs.

// Weak pointers are a great fit for scenarios like:

// Caching temporary data.
// Monitoring objects without preventing cleanup.
// Tracking objects with limited lifetimes.
// However, avoid using weak pointers in place of strong references when you need
// guaranteed access to an object. Always consider your application’s memory and
// performance requirements.

// Cache2 represents a thread-safe cache with weak pointers.
type Cache2[K comparable, V any] struct {
	mu    sync.Mutex
	items map[K]weak.Pointer[V] // Weak pointers to cached objects
}

// NewCache2 creates a new generic Cache instance.
func NewCache2[K comparable, V any]() *Cache2[K, V] {
	return &Cache2[K, V]{
		items: make(map[K]weak.Pointer[V]),
	}
}

// Get retrieves an item from the cache, if it's still alive.
func (c *Cache2[K, V]) Get(key K) (*V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Retrieve the weak pointer for the given key
	ptr, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Attempt to dereference the weak pointer
	val := ptr.Value()
	if val == nil {
		// Object has been reclaimed by the garbage collector
		delete(c.items, key)
		return nil, false
	}

	return val, true
}

// Set adds an item to the cache.
func (c *Cache2[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a weak pointer to the value
	c.items[key] = weak.Make(&value)
}

func ExecuterCache2() {
	// Create a cache with string keys and string values
	cache := NewCache2[string, string]()

	// Add an object to the cache
	data := "cached data"
	cache.Set("key1", data)

	// Retrieve it
	if val, ok := cache.Get("key1"); ok {
		fmt.Println("Cache hit:", *val)
	} else {
		fmt.Println("Cache miss")
	}

	// Simulate losing the strong reference
	data = ""
	runtime.GC() // Force garbage collection

	// Try to retrieve it again
	time.Sleep(1 * time.Second)
	
	if val, ok := cache.Get("key1"); ok {
		fmt.Println("Cache hit:", *val)
	} else {
		fmt.Println("Cache miss")
	}
}
