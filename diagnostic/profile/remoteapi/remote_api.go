package remoteapi

import (
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
)

func message(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// commands:
// go tool pprof -alloc_space http://localhost:8080/debug/pprof/heap
// go tool pprof -inuse_space http://localhost:8080/debug/pprof/heap
func Executer() {
	r := http.NewServeMux()

	r.HandleFunc("/", message)

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	http.ListenAndServe(":8080", r)
}

// allocMemory allocates 100MBs of memory
func allocMemory() []byte {
	return make([]byte, 1024*1024*100)
}

func SingleEndpoint() {
	var wg sync.WaitGroup

	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	for range 10 {
		b := allocMemory()

		log.Printf("len(b): %d", len(b))
	}

	wg.Add(1)
	wg.Wait() // this is for the benefit of the pprof server analysis
}
