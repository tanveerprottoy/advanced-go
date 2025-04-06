package statevent

import (
	"log"
	"runtime"
	"runtime/debug"
)

// The runtime provides stats and reporting of internal events for users to diagnose
// performance and utilization problems at the runtime level.
// Users can monitor these stats to better understand the overall health and
// performance of Go programs. Some frequently monitored stats and states:

// runtime.ReadMemStats reports the metrics related to heap allocation and garbage
// collection. Memory stats are useful for monitoring how much memory resources a
// process is consuming, whether the process can utilize memory well, and to catch
// memory leaks.

// debug.ReadGCStats reads statistics about garbage collection. It is useful to see how
// much of the resources are spent on GC pauses. It also reports a timeline of garbage
// collector pauses and pause time percentiles.

// debug.Stack returns the current stack trace. Stack trace is useful to see how many
// goroutines are currently running, what they are doing, and whether they are blocked or not.

// debug.WriteHeapDump suspends the execution of all goroutines and allows you to dump
// the heap to a file. A heap dump is a snapshot of a Go process' memory at a given
// time. It contains all allocated objects as well as goroutines, finalizers, and more.

// runtime.NumGoroutine returns the number of current goroutines. The value can be
// monitored to see whether enough goroutines are utilized, or to detect goroutine leaks.

func printMemStats(m *runtime.MemStats) {
	log.Printf("Alloc = %v\n", m.Alloc)
	log.Printf("TotalAlloc = %v\n", m.TotalAlloc)
	log.Printf("Sys = %v\n", m.Sys)
	log.Printf("HeapAlloc = %v\n", m.HeapAlloc)
	log.Printf("HeapSys = %v\n", m.HeapSys)
	log.Printf("HeapIdle = %v\n", m.HeapIdle)
	log.Printf("HeapInuse = %v\n", m.HeapInuse)
	log.Printf("HeapReleased = %v\n", m.HeapReleased)
	log.Printf("HeapObjects = %v\n", m.HeapObjects)
	log.Printf("StackInuse = %v\n", m.StackInuse)
	log.Printf("StackSys = %v\n", m.StackSys)
	log.Printf("NextGC = %v\n", m.NextGC)
	log.Printf("LastGC = %v\n", m.LastGC)
	log.Printf("PauseTotalNs = %v\n", m.PauseTotalNs)
	log.Printf("NumGC = %v\n", m.NumGC)
	log.Printf("GCCPUFraction = %v\n", m.GCCPUFraction)
}

// allocMemory allocates 100MBs of memory
func allocMemory() []byte {
	return make([]byte, 1024*1024*100)
}

func readMemStats() {
	// Read memory statistics before the allocation.
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	log.Println("Memory Stats Before Allocation:")
	printMemStats(&memStats)

	// Allocate some memory.
	b := allocMemory()

	log.Println("len(b): ", len(b))

	// Read memory statistics after the allocation.
	runtime.ReadMemStats(&memStats)

	log.Println("Memory Stats After Allocation:")
	printMemStats(&memStats)
}

func readGCStats() {
	var gcStats debug.GCStats

	debug.ReadGCStats(&gcStats)

	log.Printf("GC Stats: %v\n", gcStats)
}

func stack() {
	buf := make([]byte, 4096)

	// returns the current stack trace. Stack trace is useful to
	// see how many goroutines are currently running, what they
	// are doing, and whether they are blocked or not.
	// l = stack length
	l := runtime.Stack(buf, true)

	log.Printf("Stack trace: %s\n", buf[:l])
}

func heapDump() {
	// suspends the execution of all goroutines and allows you to 
	// dump the heap to a file. A heap dump is a snapshot of a Go 
	// process' memory at a given time. It contains all allocated 
	// objects as well as goroutines, finalizers, and more.
	debug.WriteHeapDump(1)
}

func numGoroutine() {
	// returns the number of current goroutines. The value can be
	// monitored to see whether enough goroutines are utilized, or 
	// to detect goroutine leaks.
	log.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())
}

func Executer() {
	readMemStats()

	readGCStats()

	stack()

	numGoroutine()
}
