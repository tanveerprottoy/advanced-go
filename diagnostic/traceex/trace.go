package traceex

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
)

func toBeTraced() {
	fmt.Printf("this function will be traced")
}

// Example demonstrates the use of the trace package to trace
// the execution of a Go program. The trace output will be
// written to the file trace.out
func Executer() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}

	defer trace.Stop()

	toBeTraced()
}
