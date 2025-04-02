package weakpointer

import (
	"runtime"
	"sync"
	"weak"
)

// This example is a little complicated, but the gist is simple.
// We start with a global concurrent map of all the mapped files we made.
// NewCachedMemoryMappedFile consults this map for an existing mapped file,
// and if that fails, creates and tries to insert a new mapped file. This could of
// course fail as well since we’re racing with other insertions, so we need to be
// careful about that too, and retry. (This design has a flaw in that we might
// wastefully map the same file multiple times in a race, and we’ll have to throw
// it away via the cleanup added by NewMemoryMappedFile. This is probably not a big
// deal most of the time. Fixing it is left as an exercise for the reader.)
// Let’s look at some useful properties of weak pointers and cleanups exploited by this code.
// Firstly, notice that weak pointers are comparable. Not only that, weak pointers have
// a stable and independent identity, which remains even after the objects they point
// to are long gone. This is why it is safe for the cleanup function to call sync.Map’s
// CompareAndDelete, which compares the weak.Pointer, and a crucial reason this code
// works at all.
// Secondly, observe that we can add multiple independent cleanups to a single
// MemoryMappedFile object. This allows us to use cleanups in a composable way and use
// them to build generic data structures. In this particular example, it might be more
// efficient to combine NewCachedMemoryMappedFile with NewMemoryMappedFile and have
// them share a cleanup. However, the advantage of the code we wrote above is that
// it can be rewritten in a generic way!
type Cache[K comparable, V any] struct {
	create func(K) (*V, error)
	m      sync.Map
}

func NewCache[K comparable, V any](create func(K) (*V, error)) *Cache[K, V] {
	return &Cache[K, V]{create: create}
}

func (c *Cache[K, V]) Get(key K) (*V, error) {
	var newValue *V
	for {
		// Try to load an existing value out of the cache.
		value, ok := cache.Load(key)
		if !ok {
			// No value found. Create a new mapped file if needed.
			if newValue == nil {
				var err error
				
				newValue, err = c.create(key)
				if err != nil {
					return nil, err
				}
			}

			// Try to install the new mapped file.
			wp := weak.Make(newValue)
			
			var loaded bool
			
			value, loaded = cache.LoadOrStore(key, wp)
			if !loaded {
				runtime.AddCleanup(newValue, func(key K) {
					// Only delete if the weak pointer is equal. If it's not, someone
					// else already deleted the entry and installed a new mapped file.
					cache.CompareAndDelete(key, wp)
				}, key)
				
				return newValue, nil
			}
		}

		// See if our cache entry is valid.
		if mf := value.(weak.Pointer[V]).Value(); mf != nil {
			return mf, nil
		}

		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		cache.CompareAndDelete(key, value)
	}
}
