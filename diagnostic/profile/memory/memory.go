package memory

import (
	"log"
	"os"
	"runtime/pprof"
)

// allocMemory allocates 100MBs of memory
func allocMemory() []byte {
	return make([]byte, 1024*1024*100)
}

func profile() {
	for range 10 {
		b := allocMemory()

		log.Printf("len(b): %d", len(b))
	}

	pprof.WriteHeapProfile(os.Stdout)
}

// commands for profiling
// go build -o app && time ./app > memory.profile
// go tool pprof memory.profile
func Executer() {
	profile()
}
