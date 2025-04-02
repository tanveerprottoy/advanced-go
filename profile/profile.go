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

func Executer() {
	flag.Parse()

	startCPUProfile()

	startMemProfile()
}
