package gccfg

import (
	"log"
	"os"
	"runtime/debug"
	"runtime/metrics"
)

// GOGC may be configured through either the GOGC environment variable
// (which all Go programs recognize), or through the SetGCPercent API in
// the runtime/debug package.
// GOGC may also be used to turn off the GC entirely (provided the memory limit does not apply)
// by setting GOGC=off or calling SetGCPercent(-1). Conceptually, this setting is equivalent
// to setting GOGC to a value of infinity, as the amount of new memory before a GC is
// triggered is unbounded.
func gc() {
	log.Printf("GOGC: %s\n", os.Getenv("GOGC"))

	debug.SetGCPercent(90)

	log.Printf("GOGC after set: %s\n", os.Getenv("GOGC"))
}

// The memory limit may be configured either via the GOMEMLIMIT environment variable which
// all Go programs recognize, or through the SetMemoryLimit function available in the
// runtime/debug package.
// This memory limit sets a maximum on the total amount of memory that the Go runtime can use.
// The specific set of memory included is defined in terms of runtime.MemStats as the expression
// Sys - HeapReleased or equivalently in terms of the runtime/metrics package,
// /memory/classes/total:bytes - /memory/classes/heap/released:bytes
// Because the Go GC has explicit control over how much heap memory it uses,
// it sets the total heap size based on this memory limit and how much other memory the
// Go runtime uses.
func mem() {
	log.Printf("GOMEMLIMIT: %s\n", os.Getenv("GOMEMLIMIT"))

	debug.SetMemoryLimit(536870912) // 512MB

	log.Printf("GOMEMLIMIT after set: %s\n", os.Getenv("GOMEMLIMIT"))

	log.Printf("metrics: %v\n", metrics.All())
}

func Executer() {
	gc()

	mem()
}
