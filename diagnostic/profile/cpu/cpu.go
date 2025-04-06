package cpu

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
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	for range 10 {
		b := allocMemory()

		log.Printf("len(b): %d", len(b))
	}
}

// commands for profiling
// go build -o app && time ./app > cpu.profile
// go tool pprof cpu.profile
func Executer() {
	profile()
}
