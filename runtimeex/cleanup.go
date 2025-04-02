//go:build unix

package runtimeex

import (
	"os"
	"runtime"
	"syscall"
)

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
