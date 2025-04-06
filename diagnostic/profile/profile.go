package profile

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

var memProfile = flag.String("memprofile", "", "write memory profile to `file`")

func startCPUProfile() {
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)

		defer pprof.StopCPUProfile()
	}
}

func startMemProfile() {
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}

		defer f.Close() // error handling omitted for example

		runtime.GC() // get up-to-date statistics

		// Lookup("allocs") creates a profile similar to go test -memprofile.
		// Alternatively, use Lookup("heap") for a profile
		// that has inuse_space as the default index.
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

func config() {
	// Block profile shows where goroutines block waiting on synchronization primitives
	// (including timer channels). Block profile is not enabled by default; use
	// runtime.SetBlockProfileRate to enable it.
	// To include every blocking event in the profile, pass rate = 1. To turn off
	// profiling entirely, pass rate <= 0.
	runtime.SetBlockProfileRate(1)

	// Mutex profile reports the lock contentions. When you think your CPU is not fully
	// utilized due to a mutex contention, use this profile. Mutex profile is not
	// enabled by default, see runtime.SetMutexProfileFraction to enable it.
	// To turn off profiling entirely, pass rate 0. To just read the current rate,
	// pass rate < 0. (For n>1 the details of sampling may change.)

	// current rate read
	r := runtime.SetMutexProfileFraction(-1)
	log.Printf("current rate: %d", r)

	// set rate
	r = runtime.SetMutexProfileFraction(1)
	log.Printf("current rate after setting value: %d", r)
}

func Executer() {
	flag.Parse()

	startCPUProfile()

	startMemProfile()
}
