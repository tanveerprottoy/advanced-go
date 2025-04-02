package runtimeex

import "runtime"

func execGC() {
	runtime.GC()
}
