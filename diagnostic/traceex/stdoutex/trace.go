package stdoutex

import (
	"log"
	"os"
	"runtime/trace"
	"sync"
)

// allocMemory allocates 100MBs of memory
func allocMemory() []byte {
	return make([]byte, 1024*1024*100)
}

func tracer() {
	trace.Start(os.Stdout)
	defer trace.Stop()

	for range 10 {
		b := allocMemory()

		log.Printf("len(b): %d", len(b))
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var result []byte
	go func() {
		result = make([]byte, 1024*1024*50)
		log.Println("done here")
		wg.Done()
	}()

	wg.Wait()
	log.Printf("%T", result)
}

// commands:
// go build -o app time ./app > app.trace
// go tool trace app.trace
func Executer() {
	tracer()
}