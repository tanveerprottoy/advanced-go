// package ex is an example
// io package usage
package ex

import (
	"io"
	"log"
	"os"
)

func countAlphabet(r io.Reader) (int, error) {
	// create a counter
	count := 0

	// create a buffer of 1024 bytes
	buff := make([]byte, 1024)

	for {
		// read from the reader
		n, err := r.Read(buff)
		// must check and process the n (number of bytes read)
		if n > 0 {
			for _, v := range buff[:n] {
				if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
					// count the alphabet
					count++
				}
			}
		}

		// check error
		if err != nil {
			// check if the error is EOF
			if err == io.EOF {
				// return the count
				return count, nil
			} else {
				// return 0 and the error
				return 0, err
			}
		}
	}
}

func ExecuterCountAlphabet() {
	f, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	// close the file
	defer f.Close()

	c, err := countAlphabet(f)
	if err != nil {
		log.Fatalf("failed to count alphabet: %v", err)
	}

	log.Printf("count: %d", c)
}
