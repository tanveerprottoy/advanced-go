package weakpointer

import (
	"os"
	"runtime"
	"sync"
	"syscall"
	"weak"
)

// Returning to our memory-mapped file example, suppose we notice that our program
// frequently maps the same files over and over, from different goroutines that are
// unaware of each other. This is fine from a memory perspective, since all these
// mappings will share physical memory, but it results in lots of unnecessary system
// calls to map and unmap the file. This is especially bad if each goroutine reads
// only a small section of each file.
// So, let’s deduplicate the mappings by filename. (Let’s assume that our program only
// reads from the mappings, and the files themselves are never modified or renamed once
// created. Such assumptions are reasonable for system font files, for example.)
// We could maintain a map from filename to memory mapping, but then it becomes unclear
// when it’s safe to remove entries from that map. We could almost use a cleanup, if it
// weren’t for the fact that the map entry itself will keep the memory-mapped file object alive.
// Weak pointers solve this problem. A weak pointer is a special kind of pointer that
// the garbage collector ignores when deciding whether an object is reachable.
// Go 1.24’s new weak pointer type, weak.Pointer, has a Value method that returns
// either a real pointer if the object is still reachable, or nil if it is not.
// If we instead maintain a map that only weakly points to the memory-mapped file,
// we can clean up the map entry when nobody’s using it anymore! Let’s see what this
// looks like.
var cache sync.Map // map[string]weak.Pointer[MemoryMappedFile]

type MemoryMappedFile struct {
	data []byte
}

func NewMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	// Get the file's info; we need its size.
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Extract the file descriptor.
	conn, err := f.SyscallConn()
	if err != nil {
		return nil, err
	}

	var data []byte

	connErr := conn.Control(func(fd uintptr) {
		// Create a memory mapping backed by this file.
		data, err = syscall.Mmap(int(fd), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	})
	if connErr != nil {
		return nil, connErr
	}
	if err != nil {
		return nil, err
	}

	mf := &MemoryMappedFile{data: data}

	cleanup := func(data []byte) {
		syscall.Munmap(data) // ignore error
	}

	runtime.AddCleanup(mf, cleanup, data)

	return mf, nil
}

func NewCachedMemoryMappedFile(filename string) (*MemoryMappedFile, error) {
	var newFile *MemoryMappedFile
	for {
		// Try to load an existing value out of the cache.
		value, ok := cache.Load(filename)
		if !ok {
			// No value found. Create a new mapped file if needed.
			if newFile == nil {
				var err error
				newFile, err = NewMemoryMappedFile(filename)
				if err != nil {
					return nil, err
				}
			}

			// Try to install the new mapped file.
			wp := weak.Make(newFile)
			var loaded bool
			value, loaded = cache.LoadOrStore(filename, wp)
			if !loaded {
				runtime.AddCleanup(newFile, func(filename string) {
					// Only delete if the weak pointer is equal. If it's not, someone
					// else already deleted the entry and installed a new mapped file.
					cache.CompareAndDelete(filename, wp)
				}, filename)
				return newFile, nil
			}
			// Someone got to installing the file before us.
			//
			// If it's still there when we check in a moment, we'll discard newFile
			// and it'll get cleaned up by garbage collector.
		}

		// See if our cache entry is valid.
		if mf := value.(weak.Pointer[MemoryMappedFile]).Value(); mf != nil {
			return mf, nil
		}

		// Discovered a nil entry awaiting cleanup. Eagerly delete it.
		cache.CompareAndDelete(filename, value)
	}
}
